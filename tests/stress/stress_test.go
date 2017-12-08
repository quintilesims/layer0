package system

import (
	"os"
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/tests/clients"
	"github.com/quintilesims/tftest"
)

type StressTest struct {
	Terraform *tftest.TestContext
	Layer0    *clients.Layer0TestClient
}

func NewStressTest(t *testing.T, dir string, vars map[string]string) *StressTest {
	if vars == nil {
		vars = map[string]string{}
	}

	vars["endpoint"] = os.Getenv(config.ENVVAR_ENDPOINT)
	vars["token"] = os.Getenv(config.ENVVAR_TOKEN)

	tfContext := tftest.NewTestContext(t,
		tftest.Dir(dir),
		tftest.Vars(vars),
		tftest.DryRun(*dry))

	layer0 := clients.NewLayer0TestClient(t, vars["endpoint"], vars["token"])

	// initialize and download modules
	tfContext.Terraformf("get")

	return &StressTest{
		Terraform: tfContext,
		Layer0:    layer0,
	}
}
