package aws

import (
	"fmt"
	"time"

	"github.com/quintilesims/layer0/common/models"
)

func (a *AdminProvider) Logs(tail int, start, end time.Time) ([]models.LogFile, error) {
	logGroupName := a.Config.LogGroupName()
	filterPattern := fmt.Sprintf("{ $.userIdentity.sessionContext.sessionIssuer.userName = \"l0-%s-ecs-role\" }", a.Config.Instance())

	return GetLogsFromCloudTrail(a.AWS.CloudWatchLogs, logGroupName, tail, start, end, filterPattern)
}
