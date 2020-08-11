module github.com/quintilesims/layer0

go 1.14

replace github.com/Sirupsen/logrus v1.5.0 => github.com/sirupsen/logrus v1.0.6

require (
	github.com/Sirupsen/logrus v1.5.0
	github.com/aws/aws-sdk-go v1.30.4
	github.com/blang/semver v3.5.1+incompatible
	github.com/briandowns/spinner v1.10.0
	github.com/dghubble/sling v1.3.0
	github.com/docker/docker v1.13.1
	github.com/emicklei/go-restful v2.12.0+incompatible
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170926063155-7524189396c6
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052
	github.com/golang/mock v1.4.3
	github.com/guregu/dynamo v1.6.1
	github.com/hashicorp/terraform v0.11.14
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/quintilesims/sts v0.0.0-20170809211516-82d6e9731e72
	github.com/quintilesims/tftest v0.0.0-20180108221958-70597d446846
	github.com/ryanuber/columnize v2.1.0+incompatible
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli v1.22.4
	github.com/zpatrick/go-bytesize v0.0.0-20170214182126-40b68ac70b6a
	github.com/zpatrick/rclient v0.0.0-20191128000351-00dfceed505a // indirect
)
