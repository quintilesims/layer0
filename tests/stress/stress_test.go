package system

import (
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/tests/clients"
	"github.com/quintilesims/tftest"
)

func NewStressTest(b *testing.B, dir string, vars map[string]string) (*tftest.TestContext, *clients.Layer0TestClient) {
	if vars == nil {
		vars = map[string]string{}
	}

	vars["endpoint"] = config.APIEndpoint()
	vars["token"] = config.AuthToken()

	tfContext := tftest.NewTestContext(
		b,
		tftest.Dir(dir),
		tftest.Vars(vars),
		tftest.DryRun(*dry),
		tftest.Log(log),
	)

	layer0 := clients.NewLayer0TestClient(b, vars["endpoint"], vars["token"])

	// download modules using terraform get
	tfContext.Terraformf("get")

	return tfContext, layer0
}
