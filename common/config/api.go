package config

import (
	"encoding/base64"
	"fmt"
	"strings"

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

	return ValidateContext(c, requiredFlags)
}

func ParseAuthToken(c *cli.Context) (string, string, error) {
	encoded := c.String(FlagToken.GetName())
	token, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", fmt.Errorf("Auth Token is not in valid base64 format: %v", err)
	}

	split := strings.Split(string(token), ":")
	if len(split) != 2 {
		return "", "", fmt.Errorf("Auth Token must be in format 'user:pass' and base64 encoded")
	}

	return split[0], split[1], nil
}
