package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/client/mock_client"
)

func setupUnitTest(t *testing.T) (*gomock.Controller, *mock_client.MockClient, *schema.Provider) {
	ctrl := gomock.NewController(t)
	mockClient := mock_client.NewMockClient(ctrl)
	p := Provider().(*schema.Provider)

	return ctrl, mockClient, p
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatal(err)
	}
}
