package ecsbackend

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/config"
)

func TestMain(m *testing.M) {
	config.SetTestConfig()
	log.SetLevel(log.FatalLevel)
	retCode := m.Run()
	os.Exit(retCode)
}
