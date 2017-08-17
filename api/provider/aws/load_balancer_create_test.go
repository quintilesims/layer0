package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadBalancer_renderLoadBalancerRolePolicy(t *testing.T) {
	template := "{{ .Region }} {{ .AccountID }} {{ .LoadBalancerID }}"

	policy, err := renderLoadBalancerRolePolicy("region", "account_id", "lbid", template)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "region account_id lbid", policy)
}
