package config

import (
	"fmt"

	"github.com/urfave/cli"
)

func APIFlags() []cli.Flag {
	return []cli.Flag{
		FlagDebug,
		FlagToken,
		FlagInstance,
		FlagPort,
		FlagNumWorkers,
		FlagJobExpiry,
		FlagLockExpiry,
		FlagAWSAccountID,
		FlagAWSAccessKey,
		FlagAWSSecretKey,
		FlagAWSRegion,
		FlagAWSVPC,
		FlagAWSLinuxAMI,
		FlagAWSWindowsAMI,
		FlagAWSS3Bucket,
		FlagAWSInstanceProfile,
		FlagAWSJobTable,
		FlagAWSTagTable,
		FlagAWSLockTable,
		FlagAWSPublicSubnets,
		FlagAWSPrivateSubnets,
		FlagAWSLogGroup,
		FlagAWSSSHKey,
		FlagAWSRequestDelay,
	}
}

func ValidateAPIContext(c *cli.Context) error {
	requiredFlags := []cli.Flag{
		FlagToken,
		FlagInstance,
		FlagAWSAccountID,
		FlagAWSAccessKey,
		FlagAWSSecretKey,
		FlagAWSVPC,
		FlagAWSLinuxAMI,
		FlagAWSWindowsAMI,
		FlagAWSS3Bucket,
		FlagAWSInstanceProfile,
		FlagAWSJobTable,
		FlagAWSTagTable,
		FlagAWSLockTable,
		FlagAWSPublicSubnets,
		FlagAWSPrivateSubnets,
		FlagAWSLogGroup,
		FlagAWSSSHKey,
	}

	for _, flag := range requiredFlags {
		name := flag.GetName()
		if !c.IsSet(name) {
			return fmt.Errorf("Required Variable %s is not set!", name)
		}
	}

	return nil
}
