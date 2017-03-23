package alicloud

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/slb"
	"regexp"
)

// common
func validateInstancePort(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 65535 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid instance port between 1 and 65535",
			k))
		return
	}
	return
}

func validateInstanceProtocol(v interface{}, k string) (ws []string, errors []error) {
	protocal := v.(string)
	if !isProtocalValid(protocal) {
		errors = append(errors, fmt.Errorf(
			"%q is an invalid value. Valid values are either http, https, tcp or udp",
			k))
		return
	}
	return
}

// ecs
func validateDiskCategory(v interface{}, k string) (ws []string, errors []error) {
	category := ecs.DiskCategory(v.(string))
	if category != ecs.DiskCategoryCloud && category != ecs.DiskCategoryCloudEfficiency && category != ecs.DiskCategoryCloudSSD {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s %s", k, ecs.DiskCategoryCloud, ecs.DiskCategoryCloudEfficiency, ecs.DiskCategoryCloudSSD))
	}

	return
}

func validateInstanceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateInstanceDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters", k))

	}
	return
}

func validateDiskName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return
	}

	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateDiskDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters", k))

	}
	return
}

//security group
func validateSecurityGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateSecurityGroupDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters", k))

	}
	return
}

func validateSecurityRuleType(v interface{}, k string) (ws []string, errors []error) {
	rt := GroupRuleDirection(v.(string))
	if rt != GroupRuleIngress && rt != GroupRuleEgress {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s", k, GroupRuleIngress, GroupRuleEgress))
	}

	return
}

func validateSecurityRuleIpProtocol(v interface{}, k string) (ws []string, errors []error) {
	pt := GroupRuleIpProtocol(v.(string))
	if pt != GroupRuleTcp && pt != GroupRuleUdp && pt != GroupRuleIcmp && pt != GroupRuleGre && pt != GroupRuleAll {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s %s %s %s", k,
			GroupRuleTcp, GroupRuleUdp, GroupRuleIcmp, GroupRuleGre, GroupRuleAll))
	}

	return
}

func validateSecurityRuleNicType(v interface{}, k string) (ws []string, errors []error) {
	pt := GroupRuleNicType(v.(string))
	if pt != GroupRuleInternet && pt != GroupRuleIntranet {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s", k, GroupRuleInternet, GroupRuleIntranet))
	}

	return
}

func validateSecurityRulePolicy(v interface{}, k string) (ws []string, errors []error) {
	pt := GroupRulePolicy(v.(string))
	if pt != GroupRulePolicyAccept && pt != GroupRulePolicyDrop {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s", k, GroupRulePolicyAccept, GroupRulePolicyDrop))
	}

	return
}

func validateSecurityPriority(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 100 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid authorization policy priority between 1 and 100",
			k))
		return
	}
	return
}

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}

func validateRouteEntryNextHopType(v interface{}, k string) (ws []string, errors []error) {
	nht := ecs.NextHopType(v.(string))
	if nht != ecs.NextHopIntance && nht != ecs.NextHopTunnel {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s", k,
			ecs.NextHopIntance, ecs.NextHopTunnel))
	}

	return
}

func validateSwitchCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
		return
	}

	mark, _ := strconv.Atoi(strings.Split(ipnet.String(), "/")[1])
	if mark < 16 || mark > 29 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a network CIDR which mark between 16 and 29",
			k))
	}

	return
}

// validateIoOptimized ensures that the string value is a valid IoOptimized that
// represents a IoOptimized - it adds an error otherwise
func validateIoOptimized(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		ioOptimized := ecs.IoOptimized(value)
		if ioOptimized != ecs.IoOptimizedNone &&
			ioOptimized != ecs.IoOptimizedOptimized {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid IoOptimized, expected %s or %s, got %q",
				k, ecs.IoOptimizedNone, ecs.IoOptimizedOptimized, ioOptimized))
		}
	}

	return
}

// validateInstanceNetworkType ensures that the string value is a classic or vpc
func validateInstanceNetworkType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		network := InstanceNetWork(value)
		if network != ClassicNet &&
			network != VpcNet {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceNetworkType, expected %s or %s, go %q",
				k, ClassicNet, VpcNet, network))
		}
	}
	return
}

func validateInstanceChargeType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		chargeType := common.InstanceChargeType(value)
		if chargeType != common.PrePaid &&
			chargeType != common.PostPaid {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceChargeType, expected %s or %s, got %q",
				k, common.PrePaid, common.PostPaid, chargeType))
		}
	}

	return
}

func validateInternetChargeType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		chargeType := common.InternetChargeType(value)
		if chargeType != common.PayByBandwidth &&
			chargeType != common.PayByTraffic {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceChargeType, expected %s or %s, got %q",
				k, common.PayByBandwidth, common.PayByTraffic, chargeType))
		}
	}

	return
}

func validateInternetMaxBandWidthOut(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 100 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid internet bandwidth out between 1 and 1000",
			k))
		return
	}
	return
}

// SLB
func validateSlbName(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		if len(value) < 1 || len(value) > 80 {
			errors = append(errors, fmt.Errorf(
				"%q must be a valid load balancer name characters between 1 and 80",
				k))
			return
		}
	}

	return
}

func validateSlbInternetChargeType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		chargeType := common.InternetChargeType(value)

		if chargeType != "paybybandwidth" &&
			chargeType != "paybytraffic" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceChargeType, expected %s or %s, got %q",
				k, "paybybandwidth", "paybytraffic", value))
		}
	}

	return
}

func validateSlbBandwidth(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 1000 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid load balancer bandwidth between 1 and 1000",
			k))
		return
	}
	return
}

func validateSlbListenerBandwidth(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if (value < 1 || value > 1000) && value != -1 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid load balancer bandwidth between 1 and 1000 or -1",
			k))
		return
	}
	return
}

func validateSlbListenerScheduler(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		scheduler := slb.SchedulerType(value)

		if scheduler != "wrr" && scheduler != "wlc" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid SchedulerType, expected %s or %s, got %q",
				k, "wrr", "wlc", value))
		}
	}

	return
}

func validateSlbListenerStickySession(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		flag := slb.FlagType(value)

		if flag != "on" && flag != "off" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid StickySession, expected %s or %s, got %q",
				k, "on", "off", value))
		}
	}
	return
}

func validateSlbListenerStickySessionType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		flag := slb.StickySessionType(value)

		if flag != "insert" && flag != "server" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid StickySessionType, expected %s or %s, got %q",
				k, "insert", "server", value))
		}
	}
	return
}

func validateSlbListenerCookie(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		flag := slb.StickySessionType(value)

		if flag != "insert" && flag != "server" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid StickySessionType, expected %s or %s, got %q",
				k, "insert", "server", value))
		}
	}
	return
}

func validateSlbListenerPersistenceTimeout(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 86400 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid load balancer persistence timeout between 0 and 86400",
			k))
		return
	}
	return
}

//data source validate func
//data_source_alicloud_image
func validateNameRegex(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if _, err := regexp.Compile(value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid regular expression: %s",
			k, err))
	}
	return
}

func validateImageOwners(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		owners := ecs.ImageOwnerAlias(value)
		if owners != ecs.ImageOwnerSystem &&
			owners != ecs.ImageOwnerSelf &&
			owners != ecs.ImageOwnerOthers &&
			owners != ecs.ImageOwnerMarketplace &&
			owners != ecs.ImageOwnerDefault {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid Image owner , expected %s, %s, %s, %s or %s, got %q",
				k, ecs.ImageOwnerSystem, ecs.ImageOwnerSelf, ecs.ImageOwnerOthers, ecs.ImageOwnerMarketplace, ecs.ImageOwnerDefault, owners))
		}
	}
	return
}

func validateRegion(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		region := common.Region(value)
		var valid string
		for _, re := range common.ValidRegions {
			if region == re {
				return
			}
			valid = valid + ", " + string(re)
		}
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid Region ID , expected %#v, got %q",
			k, valid, value))

	}
	return
}
