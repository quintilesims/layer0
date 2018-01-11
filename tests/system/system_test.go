package system

import (
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

	vars["endpoint"] = config.APIEndpoint()
	vars["token"] = config.AuthToken()

	tfContext := tftest.NewTestContext(t,
		tftest.Dir(dir),
		tftest.Vars(vars),
		tftest.DryRun(*dry))

	layer0 := clients.NewLayer0TestClient(t, vars["endpoint"], vars["token"])

	if _, err := tfContext.Terraformf("init"); err != nil {
		t.Fatal(err)
	}

	return &SystemTest{
		Terraform: tfContext,
		Layer0:    layer0,
	}
}
