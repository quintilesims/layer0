package config

import (
	"github.com/urfave/cli"
)

// common flags
var (
	FlagDebug = cli.BoolFlag{
		Name:   "d, debug",
		EnvVar: "LAYER0_DEBUG",
		Usage:  "show debug output",
	}
	FlagEndpoint = cli.StringFlag{
		Name:   "endpoint",
		Value:  DefaultEndpoint,
		EnvVar: "LAYER0_ENDPOINT",
		Usage:  "endpoint of the Layer0 API",
	}
	FlagToken = cli.StringFlag{
		Name:   "token",
		EnvVar: "LAYER0_TOKEN",
		Usage:  "auth token of the Layer0 API",
	}
	FlagAWSAccessKey = cli.StringFlag{
		Name:   "aws-access-key",
		EnvVar: "LAYER0_AWS_ACCESS_KEY",
		Usage:  "access key portion of an AWS key",
	}
	FlagAWSSecretKey = cli.StringFlag{
		Name:   "aws-secret-key",
		EnvVar: "LAYER0_AWS_SECRET_KEY",
		Usage:  "secret key portion of an AWS key",
	}
	FlagAWSRegion = cli.StringFlag{
		Name:   "aws-region",
		Value:  DefaultAWSRegion,
		EnvVar: "LAYER0_AWS_REGION",
	}
)

// api flags
var (
	FlagInstance = cli.StringFlag{
		Name:   "instance",
		EnvVar: "LAYER0_INSTANCE",
	}
	FlagPort = cli.IntFlag{
		Name:   "p, port",
		Value:  DefaultPort,
		EnvVar: "LAYER0_PORT",
	}
	FlagNumWorkers = cli.IntFlag{
		Name:   "num-workers",
		Value:  DefaultNumWorkers,
		EnvVar: "LAYER0_NUM_WORKERS",
	}
	FlagJobExpiry = cli.DurationFlag{
		Name:   "job-expiry",
		Value:  DefaultJobExpiry,
		EnvVar: "LAYER0_JOB_EXPIRY",
	}
	FlagLockExpiry = cli.DurationFlag{
		Name:   "lock-expiry",
		Value:  DefaultLockExpiry,
		EnvVar: "LAYER0_LOCK_EXPIRY",
	}
	FlagAWSAccountID = cli.StringFlag{
		Name:   "aws-account-id",
		EnvVar: "LAYER0_AWS_ACCOUNT_ID",
	}
	FlagAWSVPC = cli.StringFlag{
		Name:   "aws-vpc",
		EnvVar: "LAYER0_AWS_VPC",
	}
	FlagAWSLinuxAMI = cli.StringFlag{
		Name:   "aws-linux-ami",
		EnvVar: "LAYER0_AWS_LINUX_AMI",
	}
	FlagAWSWindowsAMI = cli.StringFlag{
		Name:   "aws-windows-ami",
		EnvVar: "LAYER0_AWS_WINDOWS_AMI",
	}
	FlagAWSS3Bucket = cli.StringFlag{
		Name:   "aws-s3-bucket",
		EnvVar: "LAYER0_AWS_S3_BUCKET",
	}
	FlagAWSInstanceProfile = cli.StringFlag{
		Name:   "aws-instance-profile",
		EnvVar: "LAYER0_AWS_INSTANCE_PROFILE",
	}
	FlagAWSJobTable = cli.StringFlag{
		Name:   "aws-job-table",
		EnvVar: "LAYER0_AWS_JOB_TABLE",
	}
	FlagAWSTagTable = cli.StringFlag{
		Name:   "aws-tag-table",
		EnvVar: "LAYER0_AWS_TAG_TABLE",
	}
	FlagAWSLockTable = cli.StringFlag{
		Name:   "aws-lock-table",
		EnvVar: "LAYER0_AWS_LOCK_TABLE",
	}
	FlagAWSPublicSubnets = cli.StringSliceFlag{
		Name:   "aws-public-subnets",
		EnvVar: "LAYER0_AWS_PUBLIC_SUBNETS",
	}
	FlagAWSPrivateSubnets = cli.StringSliceFlag{
		Name:   "aws-private-subnets",
		EnvVar: "LAYER0_AWS_PRIVATE_SUBNETS",
	}
	FlagAWSLogGroup = cli.StringFlag{
		Name:   "aws-log-group",
		EnvVar: "LAYER0_AWS_LOG_GROUP",
	}
	FlagAWSSSHKey = cli.StringFlag{
		Name:   "aws-ssh-key",
		EnvVar: "LAYER0_AWS_SSH_KEY",
	}
	FlagAWSRequestDelay = cli.DurationFlag{
		Name:   "aws-request-delay",
		Value:  DefaultAWSRequestDelay,
		EnvVar: "LAYER0_AWS_REQUEST_DELAY",
	}
)

// cli flags
var (
	FlagOutput = cli.StringFlag{
		Name:   "o, output",
		Value:  DefaultOutput,
		EnvVar: "LAYER0_OUTPUT",
		Usage:  "output format [text,json]",
	}
	FlagTimeout = cli.DurationFlag{
		Name:   "t, timeout",
		Value:  DefaultTimeout,
		EnvVar: "LAYER0_TIMEOUT",
		Usage:  "timeout [h,m,s,ms]",
	}
	FlagSkipVerifySSL = cli.BoolFlag{
		Name:   "skip-verify-ssl",
		EnvVar: "LAYER0_SKIP_VERIFY_SSL",
		Usage:  "if set, will skip ssl verification",
	}
	FlagSkipVerifyVersion = cli.BoolFlag{
		Name:   "skip-verify-version",
		EnvVar: "LAYER0_SKIP_VERIFY_VERSION",
		Usage:  "if set, will skip version verification",
	}
	FlagNoWait = cli.BoolFlag{
		Name:   "no-wait",
		EnvVar: "LAYER0_NO_WAIT",
		Usage:  "if set, will not wait for job operations to complete",
	}
)
