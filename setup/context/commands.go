package context

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/quintilesims/layer0/common/aws/provider"
	"github.com/quintilesims/layer0/common/aws/s3"
	"github.com/quintilesims/layer0/common/config"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type TFState struct {
	Version int        `json:"version"`
	Serial  int        `json:"serial"`
	Modules []TFModule `json:"modules"`
}

type TFModule struct {
	Path      []string               `json:"path"`
	Outputs   map[string]interface{} `json:"outputs"`
	Resources map[string]interface{} `json:"resources"`
}

func Apply(c *Context, force bool, dockercfg string) error {
	if err := c.Load(false); err != nil {
		return err
	}

	var vpc_id string
	vpc_flag, ok := c.Flags["vpc_id"]
	if ok && vpc_flag != nil && *vpc_flag != "" {
		fmt.Printf("Detected vpc_id: `%s`\n", *vpc_flag)
		vpc_id = *vpc_flag
		if err := addVPCtoContext(c, vpc_id); err != nil {
			return fmt.Errorf("Failed while looking up VPC %s.", vpc_id)
		}
	}

	isFirstApply, err := isFirstApply(c)
	if err != nil {
		return err
	}

	if vpc_id != "" && !isFirstApply {
		return fmt.Errorf("`apply --vpc` cannot be run on an existing layer0.  Use a new prefix to retry.")
	}

	if isFirstApply && vpc_id == "" {
		// if we're about to make a new vpc, check if this is a re-create
		vpcExists, err := vpcExists(c)
		if err != nil {
			return err
		}

		if vpcExists {
			return fmt.Errorf("`apply` would recreate vpc %s.  Use `restore` to use an existing layer0, or --vpc to install into an existing vpc.", c.Instance)
		}
	}

	if dockercfg != "" {
		fmt.Printf("Detected provided dockercfg path: `%s`\n", dockercfg)
		instancePath := fmt.Sprintf("%s/dockercfg", c.InstanceDir)

		if err := CopyFile(dockercfg, instancePath); err != nil {
			return fmt.Errorf("Failed to copy dockercfg: %v.", err)
		}
	}

	// validate dockercfg; if file does not exist, write empty dockercfg file
	if err := checkDockercfg(c, force); err != nil {
		return err
	}

	// first apply typically fails due to https://github.com/hashicorp/terraform/issues/2349
	// adding a 2nd attempt here to counter for now
	if _, err := c.Terraformf(true, "apply"); err != nil {
		fmt.Println("[WARNING] First apply attempt failed, trying again")
		if _, err := c.Terraformf(true, "apply"); err != nil {
			return err
		}
	}

	if err := Backup(c); err != nil {
		return err
	}

	if err := waitForHealthy(c, time.Minute*15); err != nil {
		return err
	}

	if isFirstApply {
		if err := initDBWithRetries(c, 5); err != nil {
			return err
		}
	}

	text := fmt.Sprintf("Successfully Applied your Layer0 '%s'! \n", c.Instance)

	if isFirstApply {
		text += "You will need to configure your LAYER0_ environment variables to connect "
		text += "the Layer0 CLI with this Layer0 instance. \n"
		text += "Please use the command './l0-setup endpoint -h' for more information."
	}

	fmt.Println(text)
	return nil
}

func isFirstApply(c *Context) (bool, error) {
	if _, err := os.Stat(c.StateFile); err != nil {
		return true, nil
	}

	file, err := ioutil.ReadFile(c.StateFile)
	if err != nil {
		return false, err
	}

	var state TFState
	if err := json.Unmarshal(file, &state); err != nil {
		return false, err
	}

	for _, module := range state.Modules {
		if len(module.Resources) != 0 {
			return false, nil
		}
	}

	return true, nil
}

func waitForHealthy(c *Context, timeout time.Duration) error {
	url := "/admin/health"
	makeCall, err := setupAPICall(c, "GET", url, nil)
	if err != nil {
		return err
	}

	fmt.Println("[INFO] Waiting for API to be healthy")
	time.Sleep(time.Minute * 1)

	for start := time.Now(); time.Since(start) < timeout; time.Sleep(time.Second * 30) {
		resp, err := makeCall()
		if err != nil {
			fmt.Printf("[WARNING] Failed to check if API is healthy: %v\n", err)
			continue
		}

		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			fmt.Println("[INFO] API is healthy")
			return nil
		}

		fmt.Printf("[WARNING] API health returned a non-200 response: %v\n", resp.Status)
	}

	return fmt.Errorf("API isn't healthy after %v", timeout)
}

func initDBWithRetries(c *Context, retries int) error {
	url := "/admin/sql"
	data := []byte(`{ "version": "latest" }`)
	makeCall, err := setupAPICall(c, "POST", url, data)
	if err != nil {
		return err
	}

	for i := 0; i < retries; i++ {
		fmt.Printf("[INFO] Attempting to update database %d/%d\n", i+1, retries)

		resp, err := makeCall()
		if err != nil {
			fmt.Printf("[WARNING] Failed to update database: %s\n", err.Error())
			continue
		}

		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			fmt.Println("[INFO] Database has been updated")
			return nil
		}

		fmt.Printf("[WARNING] Update database call returned a non-200 response: %v\n", resp.Status)
	}

	return fmt.Errorf("Failed to initialize database after %d attempts", retries)
}

func setupAPICall(c *Context, method, url string, data []byte) (func() (*http.Response, error), error) {
	endpoint, err := getTerraformOutputVariable(c, false, "endpoint")
	if err != nil {
		return nil, err
	}

	auth_token, err := getTerraformOutputVariable(c, false, "api_auth_token")
	if err != nil {
		return nil, err
	}

	url = fmt.Sprintf("%s%s", endpoint, url)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth_token))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return func() (*http.Response, error) { return client.Do(req) }, nil
}

func Backup(c *Context) error {
	return s3Action(c, true, backupFile)
}

func backupFile(conn s3.Provider, bucket, key, file string) error {
	fmt.Printf("[INFO] Backing up %s to %s/%s\n", file, bucket, key)

	if err := conn.PutObjectFromFile(bucket, key, file); err != nil {
		return fmt.Errorf("Failed to backup %s: %s", file, err.Error())
	}

	return nil
}

func Restore(c *Context) error {
	if err := s3Action(c, false, restoreFile); err != nil {
		return err
	}

	// load the tfvars we just downloaded so any further operations
	// use the updated tfvars (e.g. context.Save() in main)
	return c.Load(true)
}

func restoreFile(conn s3.Provider, bucket, key, file string) error {
	fmt.Printf("[INFO] Restoring %s/%s to %s\n", bucket, key, file)

	if err := conn.GetObjectToFile(bucket, key, file, os.FileMode(0644)); err != nil {
		return fmt.Errorf("Failed to restore %s: %s", file, err.Error())
	}

	return nil
}

func s3Action(c *Context, mustExist bool, action func(s3.Provider, string, string, string) error) error {
	if err := c.Load(mustExist); err != nil {
		return err
	}

	accessKey, secretKey, region, err := getAWSVars(c, "s3Action")
	if err != nil {
		return err
	}

	bucket := c.TerraformVars["s3_bucket"]
	if bucket == "" {
		val, err := c.PromptTFVar("s3_bucket", "S3 Bucket")
		if err != nil {
			return err
		}

		bucket = val
	}

	creds := provider.NewExplicitCredProvider(accessKey, secretKey)
	conn, err := s3.NewS3(creds, region)
	if err != nil {
		return fmt.Errorf("Failed to connect to S3: %s", err.Error())
	}

	stateKey := "terraform/terraform.tfstate"
	if err := action(conn, bucket, stateKey, c.StateFile); err != nil {
		return err
	}

	varsKey := "terraform/terraform.tfvars"
	if err := action(conn, bucket, varsKey, c.VarsFile); err != nil {
		return err
	}

	dockerConfigKey := "bootstrap/dockercfg"
	if err := action(conn, bucket, dockerConfigKey, c.DockerConfigFile); err != nil {
		return err
	}

	return nil
}

func Destroy(c *Context, force bool) error {
	if err := c.Load(true); err != nil {
		return err
	}

	if force {
		fmt.Println("Force flag present")
	} else {
		text := "Do you really want to destroy your Layer0?\n"
		text += "    This will delete all your managed infrastructure.\n"
		text += fmt.Sprintf("    There is no undo. Only '%s' will be accepted to confirm.\n", c.Instance)
		text += "\n    Enter the name of your Layer0: "

		if !requireInput(text, c.Instance) {
			return fmt.Errorf("Destroy cancelled")
		}
	}

	if _, err := c.Terraformf(true, "destroy", "--force"); err != nil {
		return err
	}

	return nil
}

func requireInput(display, requiredInput string) bool {
	fmt.Printf(display)
	var input string
	fmt.Scanln(&input)

	if input != requiredInput {
		return false
	}

	return true
}

func Endpoint(c *Context, syntax string, insecure, dev, quiet bool) error {
	if err := c.Load(true); err != nil {
		return err
	}

	var format string
	switch syntax {
	case "bash":
		format = "export %s=\"%s\"\n"
	case "powershell":
		format = "$env:%s=\"%s\"\n"
	case "cmd":
		format = "set %s=%s\n"
	default:
		return fmt.Errorf("Syntax '%s' not recognized", syntax)
	}

	settings := map[string]string{
		"endpoint":       config.API_ENDPOINT,
		"api_auth_token": config.AUTH_TOKEN,
	}

	if dev {
		settings["account_id"] = config.AWS_ACCOUNT_ID
		settings["key_pair"] = config.AWS_KEY_PAIR
		settings["agent_security_group_id"] = config.AWS_ECS_AGENT_SECURITY_GROUP_ID
		settings["ecs_instance_profile"] = config.AWS_ECS_INSTANCE_PROFILE
		settings["ecs_role"] = config.AWS_ECS_ROLE
		settings["public_subnets"] = config.AWS_PUBLIC_SUBNETS
		settings["private_subnets"] = config.AWS_PRIVATE_SUBNETS
		settings["access_key"] = config.AWS_ACCESS_KEY_ID
		settings["secret_key"] = config.AWS_SECRET_ACCESS_KEY
		settings["region"] = config.AWS_REGION
		settings["l0_prefix"] = config.PREFIX
		settings["runner_docker_image_tag"] = config.RUNNER_VERSION_TAG
		settings["vpc_id"] = config.AWS_VPC_ID
		settings["s3_bucket"] = config.AWS_S3_BUCKET
		settings["service_ami"] = config.AWS_SERVICE_AMI
	}

	for tfvar, envvar := range settings {
		val, err := getTerraformOutputVariable(c, false, tfvar)
		if err != nil {
			return err
		}

		fmt.Printf(format, envvar, val)
	}

	if dev {
		fmt.Printf(format, config.DB_CONNECTION, fmt.Sprintf("layer0:nohaxplz@tcp(localhost:3306)/"))
		fmt.Printf(format, config.DB_NAME, fmt.Sprintf("layer0_%s", c.Instance))
	}

	if insecure {
		fmt.Printf(format, config.SKIP_SSL_VERIFY, "1")
	}

	if quiet {
		fmt.Printf(format, config.SKIP_VERSION_VERIFY, "1")
	}
	fmt.Println("# Run this command to configure your shell:")
	fmt.Println("# eval $(./l0-setup endpoint -i", c.Instance, ")")

	return nil
}

func Plan(c *Context, args []string) error {
	if err := checkDockercfg(c, false); err != nil {
		return err
	}

	args = append([]string{"plan"}, args...)
	return Terraform(c, args)
}

func checkDockercfg(c *Context, force bool) error {
	path := fmt.Sprintf("%s/dockercfg", c.InstanceDir)

	// if file does not exist
	if _, err := os.Stat(path); err != nil {

		// write empty dockercfg file
		if err := ioutil.WriteFile(path, []byte("{}"), 0660); err != nil {
			return fmt.Errorf("Failed to write 'dockercfg': %v", err)
		}

		// alert user
		if !force {
			text := fmt.Sprintf("NOTICE: No 'dockercfg' file present at %s. \n", path)
			text += "Created default 'dockercfg' file. \n"
			text += "\n"
			text += "If you require private registry authentication, please edit this file \n"
			text += "with the required credentials before proceeding. \n"
			text += "\n    Continue? (y/n): "

			if !requireInput(text, "y") {
				return fmt.Errorf("Operation Cancelled")
			}
		}

	}

	return validateDockercfg(path)
}

func validateDockercfg(dockercfgPath string) error {
	contents, err := ioutil.ReadFile(dockercfgPath)
	if err != nil {
		return err
	}

	// Simple attempt to validate JSON
	var result map[string]interface{}
	if err := json.Unmarshal(contents, &result); err != nil {
		return fmt.Errorf("Failed to validate JSON for 'dockercfg': %v", err)
	}

	return nil
}

func Terraform(c *Context, args []string) error {
	if err := c.Load(false); err != nil {
		return err
	}

	if _, err := c.Terraformf(true, args...); err != nil {
		return err
	}

	return nil
}

func getTerraformOutputVariable(c *Context, showOutput bool, variable string) (string, error) {
	val, err := c.Terraformf(showOutput, "output", variable)
	if err != nil {
		return "", fmt.Errorf("Failed to get '%s' from terraform outputs: %s", variable, err.Error())
	}

	return strings.Replace(val, "\n", "", 1), nil
}
