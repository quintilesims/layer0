package system

import (
	"log"
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

	log.Printf("endpoint: %s", vars["endpoint"])
	log.Printf("token: %s", vars["token"])

	tfContext := tftest.NewTestContext(t,
		tftest.Dir(dir),
		tftest.Vars(vars),
		tftest.DryRun(*dry))

	layer0 := clients.NewLayer0TestClient(t, vars["endpoint"], vars["token"])

	// download modules using terraform init
	stdoutStderr, err := tfContext.Init()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	log.Printf("%s\n", stdoutStderr)

	return &SystemTest{
		Terraform: tfContext,
		Layer0:    layer0,
	}
}
