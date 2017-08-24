package system

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/tests/clients"
	"github.com/quintilesims/tftest"
)

type StressTest struct {
	Terraform *tftest.TestContext
	Layer0    *clients.Layer0TestClient
}

func NewStressTest(t tftest.Tester, dir string, vars map[string]string) *StressTest {
	if vars == nil {
		vars = map[string]string{}
	}

	vars["endpoint"] = config.APIEndpoint()
	vars["token"] = config.AuthToken()

	tfContext := tftest.NewTestContext(t,
		tftest.Dir(dir),
		tftest.Vars(vars),
		tftest.DryRun(*dry),
		tftest.Log(t.(tftest.Logger)))

	layer0 := clients.NewLayer0TestClient(t, vars["endpoint"], vars["token"])

	// download modules using terraform get
	tfContext.Terraformf("get")

	return &StressTest{
		Terraform: tfContext,
		Layer0:    layer0,
	}
}
