package instance

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
	"github.com/quintilesims/layer0/setup/terraform"
)

func (l *LocalInstance) Upgrade(version string, force bool) error {
	if err := l.assertExists(); err != nil {
		return err
	}

	// load terraform config from ~/.layer0/<instance>/main.tf.json
	config, err := l.loadLayer0Config()
	if err != nil {
		return err
	}

	// create the 'layer0' module if it doesn't already exist
	if _, ok := config.Modules["layer0"]; !ok {
		config.Modules["layer0"] = terraform.Module{}
	}

	module := config.Modules["layer0"]

	// only patch upgrades are allowed
	if current, ok := module["version"]; ok && !force {
		if err := assertPatchUpgrade(current.(string), version); err != nil {
			return err
		}
	}

	// set new input values for 'source' and 'version'
	inputValues := map[string]string{
		INPUT_SOURCE:  fmt.Sprintf("%s?ref=%v", LAYER0_MODULE_SOURCE, version),
		INPUT_VERSION: version,
	}

	for input, value := range inputValues {
		if current, ok := module[input]; ok && current != value && !force {
			fmt.Printf("This will update the '%s' input \n\tFrom: [%s] \n\tTo:   [%s]\n\n", input, current, value)
			fmt.Printf("\tPress 'enter' to accept this change: ")

			var input string
			fmt.Scanln(&input)
		}

		module[input] = value
	}

	// save the terraform config as ~/.layer0/<instance>/main.tf.json
	path := fmt.Sprintf("%s/main.tf.json", l.Dir)
	if err := terraform.WriteConfig(path, config); err != nil {
		return err
	}

	// run `terraform get` to download terraform modules
	if err := l.Terraform.Get(l.Dir); err != nil {
		return err
	}

	// validate the terraform configuration
	if err := l.Terraform.Validate(l.Dir); err != nil {
		return err
	}

	return nil
}

func assertPatchUpgrade(currentVersion, desiredVersion string) error {
	currentVersion = strings.TrimPrefix(currentVersion, "v")
	desiredVersion = strings.TrimPrefix(desiredVersion, "v")

	current, err := semver.Make(currentVersion)
	if err != nil {
		text := fmt.Sprintf("Failed to parse current version ('%s'): %v\n", currentVersion, err)
		text += "Use --force to disable semantic version checking"
		return fmt.Errorf(text)
	}

	desired, err := semver.Make(desiredVersion)
	if err != nil {
		return fmt.Errorf("Failed to parse desired version: %v", err)
	}

	if current.Major != desired.Major {
		return fmt.Errorf("Cannot change Major versions (current: %s)", current.String())
	}

	if current.Minor != desired.Minor {
		return fmt.Errorf("Cannot change Minor versions (current: %s)", current.String())
	}

	return nil
}
