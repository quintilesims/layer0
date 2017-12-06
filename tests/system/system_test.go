package system

import (
	"os"
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/tests/clients"
	"github.com/quintilesims/tftest"
)

type SystemTest struct {
	Terraform *tftest.TestContext
	Layer0    *clients.Layer0TestClient
}

func NewSystemTest(t *testing.T, dir string, vars map[string]string) *SystemTest {
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
	tfContext.Terraformf("init")
	tfContext.Terraformf("get")

	return &SystemTest{
		Terraform: tfContext,
		Layer0:    layer0,
	}
}
