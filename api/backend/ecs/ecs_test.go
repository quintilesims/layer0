package ecsbackend

import (
	log "github.com/Sirupsen/logrus"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.SetTestConfig()
	log.SetLevel(log.FatalLevel)
	retCode := m.Run()
	os.Exit(retCode)
}
