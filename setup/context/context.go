package context

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"text/template"
)

type Context struct {
	Instance      string
	StateFile     string
	VarsFile      string
	InstanceDir   string
	ExecutionDir  string
	Flags         map[string]*string
	TerraformVars map[string]string
	tfvarsCache   map[string]string
}

func NewContext(instance, version string, flags map[string]*string) (*Context, error) {
	if flags == nil {
		flags = map[string]*string{}
	}

	// go-homedir is required for cross-compilation
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	instanceDir := fmt.Sprintf("%s/layer0/instances/%s", homeDir, instance)

	executionDir, err := GetExecutionDir()
	if err != nil {
		fmt.Errorf("Couldn't find CWD: %v", err)
	}

	context := &Context{
		Instance:     instance,
		Flags:        flags,
		InstanceDir:  instanceDir,
		ExecutionDir: executionDir,
		StateFile:    fmt.Sprintf("%s/terraform.tfstate", instanceDir),
		VarsFile:     fmt.Sprintf("%s/terraform.tfvars", instanceDir),
		TerraformVars: map[string]string{
			"setup_version": version,
		},
	}

	return context, nil
}

func (this *Context) Load(mustExist bool) error {
	re := regexp.MustCompile("^[a-z][a-z0-9]{0,15}$")
	if !re.MatchString(this.Instance) {
		text := "INSTANCE argument violated one or more of the following constraints: \n"
		text += "1. Start with a lowercase letter \n"
		text += "2. Only contain lowercase alphanumeric characters \n"
		text += "3. Be between 1 and 15 characters in length \n"
		return fmt.Errorf(text)
	}

	if err := os.MkdirAll(this.InstanceDir, 0700); err != nil {
		return fmt.Errorf("Failed to create instance directory: %v", err)
	}

	// todo: remove all existing *.tf and *.template files

	if _, err := os.Stat(this.StateFile); err != nil && mustExist {
		text := fmt.Sprintf(" '%s' doesn't exist locally. ", this.StateFile)
		text += fmt.Sprintf("Have you tried running `l0-setup restore %s' ?", this.Instance)
		return fmt.Errorf(text)
	}

	if _, err := os.Stat(this.VarsFile); err != nil && mustExist {
		text := fmt.Sprintf(" '%s' doesn't exist locally. ", this.VarsFile)
		text += fmt.Sprintf("Have you tried running `l0-setup restore %s' ?", this.Instance)
		return fmt.Errorf(text)
	}

	tfvars, err := this.loadTFVars()
	if err != nil {
		return err
	}

	for key, val := range tfvars {
		this.TerraformVars[key] = val
	}

	for _, variable := range TerraformVariables {
		var val string
		var err error

		val, err = this.tryFindTFVar(variable.Name)
		if err != nil {
			return err
		}

		if val == "" || variable.AlwaysCalculate {
			if variable.IsCalculated {
				val, err = variable.Calculate(this)
				if err != nil {
					return err
				}
			} else if variable.Default != "" {
				val = variable.Default
			} else {
				val, err = this.PromptTFVar(variable.Name, variable.DisplayName)
				if err != nil {
					return err
				}
			}
		}

		if variable.Validate != nil {
			if err := variable.Validate(this, val); err != nil {
				return err
			}
		}

		this.TerraformVars[variable.Name] = val
	}

	return this.Save()
}

func (this *Context) tryFindTFVar(key string) (string, error) {
	// check if it was passed via command line
	if val, ok := this.Flags[key]; ok && *val != "" {
		return *val, nil
	}

	// check if it was passed using terraform-style flags: "-var key=val" or "--var key=val"
	args := os.Args
	for i, arg := range args {
		if (arg == "-var" || arg == "--var") && i+1 < len(args) {
			split := strings.Split(args[i+1], "=")
			if split[0] == key {
				return split[1], nil
			}
		}
	}

	// check TF_VAR_* environment variables for override value
	if val := os.Getenv(fmt.Sprintf("TF_VAR_%s", key)); val != "" {
		return val, nil
	}

	// check tfvars file for value
	tfvars, err := this.loadTFVars()
	if err != nil {
		return "", err
	}

	if val, ok := tfvars[key]; ok {
		return val, nil
	}

	return "", nil
}

func (this *Context) loadTFVars() (map[string]string, error) {
	tfvars := map[string]string{}

	if _, err := os.Stat(this.VarsFile); err == nil {
		lines, err := ReadFileLines(this.VarsFile)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Failed to read %s: %s", this.VarsFile, err.Error())
		}

		for _, line := range lines {
			line = strings.Replace(line, "=", "", 1)
			split := strings.Split(line, "\"")
			key := strings.Replace(split[0], " ", "", -1)
			val := split[1]
			tfvars[key] = val
		}
	}

	return tfvars, nil
}

func (this *Context) PromptTFVar(key, display string) (string, error) {
	val := ""
	for attempts := 0; val == ""; attempts++ {
		if attempts >= 5 {
			return "", fmt.Errorf("Failed to get input for variable '%s'.", key)
		}

		fmt.Printf("%s: ", display)
		fmt.Scanln(&val)
	}

	return val, nil
}

func (this *Context) Save() error {
	file, err := os.OpenFile(this.VarsFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %v", this.VarsFile, err)
	}
	defer file.Close()

	// sort keys for consistent file diff
	keys := make([]string, 0, len(this.TerraformVars))
	for key, _ := range this.TerraformVars {
		keys = append(keys, key)
	}

	sort.Sort(ByString(keys))

	for _, key := range keys {
		val := this.TerraformVars[key]
		if _, err := file.WriteString(fmt.Sprintf("%s = \"%v\"\n", key, val)); err != nil {
			return fmt.Errorf("Failed to write %s: %v", this.VarsFile, err)
		}
	}

	return nil
}

func (this *Context) Terraformf(showOutput bool, args ...string) (string, error) {
	var stdout string
	var stderr string

	readOut := func(scanner *bufio.Scanner) {
		text := scanner.Text()
		stdout += fmt.Sprintf("%s\n", text)

		if showOutput {
			fmt.Fprintln(os.Stdout, text)
		}
	}

	readErr := func(scanner *bufio.Scanner) {
		text := scanner.Text()
		stderr += fmt.Sprintf("%s\n", text)

		if showOutput {
			fmt.Fprintln(os.Stderr, text)
		}
	}

	if err := this.executeTerraform(readOut, readErr, args...); err != nil {
		return "", fmt.Errorf("%v %s", err, stderr)
	}

	return stdout, nil
}

func (this *Context) executeTerraform(readOut func(*bufio.Scanner), readErr func(*bufio.Scanner), args ...string) error {
	if err := this.addBinPath(); err != nil {
		return err
	}

	if err := this.writeTerraformFiles(); err != nil {
		return err
	}

	args = append(args, "-no-color")
	cmd := exec.Command("terraform", args...)
	cmd.Dir = this.InstanceDir

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	outScanner := bufio.NewScanner(outReader)
	go func() {
		for outScanner.Scan() {
			readOut(outScanner)
		}
	}()

	errReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	errScanner := bufio.NewScanner(errReader)
	go func() {
		for errScanner.Scan() {
			readErr(errScanner)
		}
	}()

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (this *Context) addBinPath() error {
	var delim string
	if runtime.GOOS == "windows" {
		delim = ";"
	} else {
		delim = ":"
	}

	path := fmt.Sprintf("%s/bin%s%s", this.ExecutionDir, delim, os.Getenv("PATH"))
	if err := os.Setenv("PATH", path); err != nil {
		return fmt.Errorf("Couldn't update PATH: %v", err)
	}

	return nil
}

func (this *Context) writeTerraformFiles() error {
	templates := []string{
		"autoscaling.tf.template",
		"ec2.tf.template",
		"ecs.tf.template",
		"elb.tf.template",
		"iam.tf.template",
		"outputs.tf.template",
		"provider.tf.template",
		"rds.tf.template",
		"s3.tf.template",
		"variables.tf.template",
		"vpc.tf.template",
		"cloudwatch.tf.template",
		"certificate.tf.template",
	}

	for _, fileName := range templates {
		filePath := fmt.Sprintf("%s/%s", this.ExecutionDir, fileName)
		tmpl, err := template.ParseFiles(filePath)
		if err != nil {
			return fmt.Errorf("Failed to parse template: %v", err)
		}

		outPath := strings.Replace(fileName, ".template", "", 1)
		outPath = fmt.Sprintf("%s/%s", this.InstanceDir, outPath)

		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(outFile, this.TerraformVars); err != nil {
			return fmt.Errorf("Failed to write template:  %v", err)
		}
	}

	templatesDir := fmt.Sprintf("%s/templates", this.ExecutionDir)
	if err := CopyDir(templatesDir, fmt.Sprintf("%s/templates", this.InstanceDir)); err != nil {
		return err
	}

	return nil
}
