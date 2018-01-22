package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/urfave/cli"
)

func TestCreateLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName:  "lb_name",
		InstanceType:     "t2.small",
		MinScale:         2,
		MaxScale:         5,
		UserDataTemplate: []byte("user_data"),
		OperatingSystem:  "linux",
		AMIID:            "ami",
	}

	base.Client.EXPECT().
		CreateLoadBalancer(req).
		Return("lb_id", nil)

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(&models.LoadBalancer{}, nil)

	flags := map[string]interface{}{
		"type":      req.InstanceType,
		"min-scale": req.MinScale,
		"max-scale": req.MaxScale,
		"user-data": file.Name(),
		"os":        req.OperatingSystem,
		"ami":       req.AMIID,
	}

	c := testutils.NewTestContext(t, []string{"lb_name"}, flags)
	if err := command.create(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateLoadBalancerInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
		"Negative MinScale": testutils.NewTestContext(t,
			[]string{"lb_name"},
			map[string]interface{}{"min-scale": "-1"}),
		"Negative MaxScale": testutils.NewTestContext(t,
			[]string{"lb_name"},
			map[string]interface{}{"max-scale": "-1"}),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.create(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("loadBalancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	base.Client.EXPECT().
		DeleteLoadBalancer("lb_id").
		Return(nil)

	c := testutils.NewTestContext(t, []string{"lb_name"}, nil)
	if err := command.delete(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteLoadBalancerInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.delete(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}
}

func TestListLoadBalancers(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	base.Client.EXPECT().
		ListLoadBalancers().
		Return([]models.LoadBalancerSummary{}, nil)

	c := testutils.NewTestContext(t, nil, nil)
	if err := command.list(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadLoadBalancer(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	base.Resolver.EXPECT().
		Resolve("loadBalancer", "lb_name").
		Return([]string{"lb_id"}, nil)

	base.Client.EXPECT().
		ReadLoadBalancer("lb_id").
		Return(&models.LoadBalancer{}, nil)

	c := testutils.NewTestContext(t, []string{"lb_name"}, nil)
	if err := command.read(c); err != nil {
		t.Fatal(err)
	}
}

func TestReadLoadBalancerInputErrors(t *testing.T) {
	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	command := NewLoadBalancerCommand(base.Command())

	contexts := map[string]*cli.Context{
		"Missing NAME arg": testutils.NewTestContext(t, nil, nil),
	}

	for name, c := range contexts {
		t.Run(name, func(t *testing.T) {
			if err := command.read(c); err == nil {
				t.Fatal("error was nil!")
			}
		})
	}

}
