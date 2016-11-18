package context

// utilities to look up existing vpc details

import (
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ec2"
	"gitlab.imshealth.com/xfra/layer0/common/aws/provider"
	"net"
	"sort"
	"strings"
)

func Vpc(c *Context, vpcId string) error {
	if err := c.Load(true); err != nil {
		return err
	}

	fmt.Println("VPC test - detecting if the given VPC is acceptable by Layer0")

	if err := addVPCtoContext(c, vpcId); err != nil {
		fmt.Println("Failed to lookup VPC")
		return err
	}

	fmt.Printf("VPC looks good! run `l0-setup apply --vpc %s %s` to setup.\n", vpcId, c.Instance)

	return nil
}

func vpcExists(c *Context) (bool, error) {
	// this name is generated in vpc.tf.template with this format
	vpcName := fmt.Sprintf("l0-%s-vpc", c.Instance)

	fmt.Printf("Checking vpc: %s ... ", vpcName)

	conn, err := newEC2(c)
	if err != nil {
		return false, fmt.Errorf("Failed to create EC2 connection: %s", err.Error())
	}

	vpc, err := conn.DescribeVPCByName(vpcName)
	if err != nil {
		return false, err
	}

	if vpc == nil {
		fmt.Println("No vpc returned")
		return false, nil
	}

	fmt.Printf("Found vpc: %s, %s\n", *vpc.VpcId, *vpc.CidrBlock)
	return true, nil
}

func addVPCtoContext(c *Context, vpcId string) error {
	conn, err := newEC2(c)
	if err != nil {
		return fmt.Errorf("Failed to create EC2 connection: %s", err.Error())
	}

	// match vpc
	vpcs, err := conn.DescribeVPC(vpcId)
	if err != nil {
		fmt.Println("No VPC found")
		return err
	}
	fmt.Printf("Found vpc: %s, %s\n", *vpcs.VpcId, *vpcs.CidrBlock)

	// match cidr
	cidr, prefix, err := getCIDR(vpcs.CidrBlock)
	if err != nil {
		fmt.Println("Failed to Get CIDR for VPC")
		return err
	}
	fmt.Printf("Found Cidr prefix %s\n", prefix)

	subnetCidrs, err := suggestSubnets(conn, vpcId, cidr)
	if err != nil {
		fmt.Println("Failed to create CIDR ranges for new subnets")
		return err
	}

	fmt.Printf("New Subnets: %v\n", subnetCidrs)

	// with vpc, cidr, and subnets - we can use the VPC.
	c.TerraformVars["vpc_id"] = vpcId
	c.TerraformVars["cidr_prefix"] = prefix
	c.TerraformVars["public_subnet_cidr_a"] = subnetCidrs[0]
	c.TerraformVars["public_subnet_cidr_b"] = subnetCidrs[1]
	c.TerraformVars["private_subnet_cidr_a"] = subnetCidrs[2]
	c.TerraformVars["private_subnet_cidr_b"] = subnetCidrs[3]

	// resource we (optionally) could reuse. eg: gatway, nat
	getRoutes(c, conn, vpcId, cidr.String())

	return nil
}

func newEC2(c *Context) (ec2.Provider, error) {
	accessKey, secretKey, region, err := getAWSVars(c, "vpc")
	if err != nil {
		return nil, err
	}

	fmt.Printf("Listing VPC in %s region\n", region)
	creds := provider.NewExplicitCredProvider(accessKey, secretKey)
	return ec2.NewEC2(creds, region)
}

func getCIDR(cidrBlock *string) (*net.IPNet, string, error) {
	if cidrBlock == nil {
		return nil, "", fmt.Errorf("vpc had not CIDR information")
	}
	vpc_cidr := *cidrBlock

	_, ipnet, err := net.ParseCIDR(vpc_cidr)
	if err != nil {
		return nil, "", fmt.Errorf("could not parse cidr: %v", err)
	}

	ones, _ := ipnet.Mask.Size()
	if ones == 16 {
		pieces := strings.Split(vpc_cidr, ".")[:2]
		prefix := strings.Join(pieces, ".")
		return ipnet, prefix, nil
	} else {
		// TODO - overly conservative.
		// as long as we can make new /24s this range is ok.
		// one way to do that might be:
		// - parse the existing subnets
		// - take the highest /24 range
		// - add ip[2] += 1, and check net.Contains(result)
		return nil, "", fmt.Errorf("prefix was not /16.  For simplicity setup can only use networks of this form")
	}
}

func suggestSubnets(conn ec2.Provider, vpcId string, cidr *net.IPNet) ([]string, error) {
	subnetCidrs := make([]*net.IPNet, 0, 2)
	nets, err := conn.DescribeVPCSubnets(vpcId)
	if err != nil {
		fmt.Printf("Error finding VPCs")
		return nil, err
	}

	// enforce that all existing subnets are /24
	for _, subnet := range nets {
		_, ipnet, err := net.ParseCIDR(*subnet.CidrBlock)
		if err == nil {
			ones, _ := ipnet.Mask.Size()
			if ones == 24 {
				subnetCidrs = append(subnetCidrs, ipnet)
			} else {
				return nil, fmt.Errorf("subnet (%s) was not /24.", *subnet.SubnetId)
			}
		} else {
			fmt.Printf("Failed to parse cidr %s -- %v\n", *subnet.CidrBlock, err)
			return nil, err
		}
	}

	// sort so we know the highest found
	sort.Sort(ByIPDesc(subnetCidrs))
	fmt.Printf("Sorted: %v\n", subnetCidrs)

	// record the highest IP subnet
	var ip_block_2 byte
	if len(subnetCidrs) > 0 {
		ip_block_2 = subnetCidrs[0].IP[2]
	} else {
		// this handles the case when the VPC has no subnets
		ip_block_2 = cidr.IP[2]
	}

	// try to use the next 4 /24 subnets after the highest discovered
	new_subnets := make([]string, 4, 4)
	for i := 0; i < 4; i++ {
		new_ip := net.IPNet{
			IP:   net.IPv4(cidr.IP[0], cidr.IP[1], ip_block_2+1+byte(i), 0),
			Mask: net.CIDRMask(24, 32),
		}

		if !cidr.Contains(new_ip.IP) {
			return nil, fmt.Errorf("The VPC CIDR block is full (generated CIDR %v was out of range)", new_ip)
		}

		new_subnets[i] = new_ip.String()
	}

	return new_subnets, nil
}

// sort functions for net.IPNet (cidr blocks)
type ByIPDesc []*net.IPNet

func (ip ByIPDesc) Len() int {
	return len(ip)
}

func (ip ByIPDesc) Swap(i, j int) {
	ip[i], ip[j] = ip[j], ip[i]
}

func (ip ByIPDesc) Less(i, j int) bool {
	left, right := ip[i].IP, ip[j].IP
	for index, block := range left {
		if block == right[index] {
			continue
		}
		return block > right[index]
	}

	return false
}

func getRoutes(c *Context, conn ec2.Provider, vpcId string, vpc_cidr string) {
	var gatewayId *string
	var natId *string
	rts, err := conn.DescribeVPCRoutes(vpcId)
	if err != nil {
		fmt.Printf("Failed describing routes")
		return
	}

	for _, rt := range rts {
		for _, route := range rt.Routes {
			if *route.State != "active" {
				// skip inactive routes
				continue
			}

			if route.InstanceId != nil {
				fmt.Printf("Private, with nat: %s\n", *route.InstanceId)
				natId = route.InstanceId

				nat_ok := isAcceptableNAT(conn, natId, vpc_cidr)
				if nat_ok {
					c.TerraformVars["nat_id"] = *natId
					c.TerraformVars["private_route_table_id"] = *rt.RouteTableId
				}
			}

			if route.GatewayId != nil {
				if *route.GatewayId == "local" {
					// skip local route, it doesn't tell us public/private
					continue
				} else {
					gatewayId = route.GatewayId
					fmt.Printf("Public with igw: %s\n", *gatewayId)
					c.TerraformVars["igw_id"] = *gatewayId
					c.TerraformVars["public_route_table_id"] = *rt.RouteTableId
				}
			}
		}

		// rt.Associations is interesting if we re-use subnets
	}
}

func isAcceptableNAT(conn ec2.Provider, natId *string, vpc_cidr string) bool {
	if natId == nil {
		return false
	}

	instance, err := conn.DescribeInstance(*natId)
	if err != nil {
		fmt.Printf("Failed describe NAT instance")
		return false
	}

	sgs := make([]string, len(instance.SecurityGroups), len(instance.SecurityGroups))
	for i, sg := range instance.SecurityGroups {
		sgs[i] = *sg.GroupName
	}

	for _, name := range sgs {
		sg, err := conn.DescribeSecurityGroup(name)
		if err != nil {
			fmt.Printf("Failed to describe securityGroup %s", name)
			return false
		}

		if isNatHttpSecurityGroup(sg, vpc_cidr) {
			fmt.Printf("Nat contains required permissions: %s\n", *sg.GroupId)
			return true
		} else {
			fmt.Printf("Security Group %s does not match requirements\n", *sg.GroupId)
		}
	}

	return false
}

func isNatHttpSecurityGroup(sg *ec2.SecurityGroup, vpc_cidr string) bool {
	// Wish-List: What the method is doing is comparing the current state of the NAT to
	// what terraform would create... terraform can do the comparison (plan) but not
	// the importing of the existing config.  So to work-around we represent the
	// terraform state in code, and manually sync with vpc.tf.template

	// check if the security group is valid for outgoing http
	foundHttp := false
	foundHttps := false
	foundIcmp := false
	for _, perm := range sg.IpPermissionsEgress {
		// for simplicity, only accept permissions that encompass the whole vpc
		foundCidr := false
		for _, ips := range perm.IpRanges {
			if ips.CidrIp != nil {
				if *ips.CidrIp == vpc_cidr {
					foundCidr = true
				} else if *ips.CidrIp == "0.0.0.0/0" {
					foundCidr = true
				}
			}
		}

		if !foundCidr {
			fmt.Println("Cidr not found")
			continue
		}

		if perm.FromPort != nil && perm.ToPort != nil {
			if *perm.FromPort <= 80 && *perm.ToPort >= 80 {
				foundHttp = true
			}
			if *perm.FromPort <= 443 && *perm.ToPort >= 443 {
				foundHttps = true
			}
		}

		if perm.IpProtocol != nil && *perm.IpProtocol == "icmp" {
			foundIcmp = true
		}
	}

	return foundHttps && foundHttp && foundIcmp
}
