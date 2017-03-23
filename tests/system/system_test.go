package system

import (
	"fmt"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/tests/system/clients"
	"github.com/quintilesims/tftest"
	"testing"
)

type SystemTest struct {
	Terraform *tftest.TestContext
	L0Client  *clients.Layer0TestClient
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
		tftest.DryRun(*dry),
		tftest.Log(log))

	// download modules using terraform get
	tfContext.Terraformf("get")

	l0Client := clients.NewLayer0TestClient(t,
		vars["endpoint"],
		fmt.Sprintf("Basic %s", vars["token"]))

	return &SystemTest{
		Terraform: tfContext,
		L0Client:  l0Client,
	}
}
