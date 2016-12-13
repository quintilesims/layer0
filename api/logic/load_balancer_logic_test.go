package logic

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/commmon/db"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestGetLoadBalancer(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.GetLoadBalancer with correct param",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					GetLoadBalancer("lbid").
					Return(&models.LoadBalancer{}, nil)

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)
				logic.GetLoadBalancer("lbid")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.GetLoadBalancer error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					GetLoadBalancer(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if _, err := logic.GetLoadBalancer("lbid"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate model with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetLoadBalancer(gomock.Any()).
					Return(&models.LoadBalancer{LoadBalancerID: "lbid"}, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "lbid",
					EntityType: "load_balancer",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "lbid",
					EntityType: "load_balancer",
					Key:        "environment_id",
					Value:      "envid",
				})

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				loadBalancer, err := logic.GetLoadBalancer("lbid")
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(loadBalancer.LoadBalancerID, "lbid")
				reporter.AssertEqual(loadBalancer.LoadBalancerName, "some_name")
				reporter.AssertEqual(loadBalancer.EnvironmentID, "envid")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetLoadBalancer(gomock.Any()).
					Return(&models.LoadBalancer{}, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if _, err := logic.GetLoadBalancer("lbid"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListLoadBalancers(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.ListLoadBalancers",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListLoadBalancers().
					Return([]*models.LoadBalancer{}, nil)

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)
				logic.ListLoadBalancers()
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.ListLoadBalancers error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListLoadBalancers().
					Return(nil, fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if _, err := logic.ListLoadBalancers(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate models with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				loadBalancers := []*models.LoadBalancer{
					&models.LoadBalancer{LoadBalancerID: "some_id_1"},
					&models.LoadBalancer{LoadBalancerID: "some_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListLoadBalancers().
					Return(loadBalancers, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "some_id_1",
					EntityType: "load_balancer",
					Key:        "name",
					Value:      "some_name_1",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "some_id_1",
					EntityType: "load_balancer",
					Key:        "environment_id",
					Value:      "some_env_1",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "some_id_2",
					EntityType: "load_balancer",
					Key:        "name",
					Value:      "some_name_2",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "some_id_2",
					EntityType: "load_balancer",
					Key:        "environment_id",
					Value:      "some_env_2",
				})

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				loadBalancers, err := logic.ListLoadBalancers()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(loadBalancers[0].LoadBalancerID, "some_id_1")
				reporter.AssertEqual(loadBalancers[0].LoadBalancerName, "some_name_1")
				reporter.AssertEqual(loadBalancers[0].EnvironmentID, "some_env_1")

				reporter.AssertEqual(loadBalancers[1].LoadBalancerID, "some_id_2")
				reporter.AssertEqual(loadBalancers[1].LoadBalancerName, "some_name_2")
				reporter.AssertEqual(loadBalancers[1].EnvironmentID, "some_env_2")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				loadBalancers := []*models.LoadBalancer{
					&models.LoadBalancer{LoadBalancerID: "some_id_1"},
					&models.LoadBalancer{LoadBalancerID: "some_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListLoadBalancers().
					Return(loadBalancers, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if _, err := logic.ListLoadBalancers(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteLoadBalancer(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.DeleteLoadBalancer with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					DeleteLoadBalancer("lbid").
					Return(nil)

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)
				logic.DeleteLoadBalancer("lbid")
			},
		},
		testutils.TestCase{
			Name: "Should delete loadBalancer tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteLoadBalancer(gomock.Any()).
					Return(nil)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "lbid",
					EntityType: "load_balancer",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "not_lbid",
					EntityType: "load_balancer",
					Key:        "name",
					Value:      "some_name",
				})

				return map[string]interface{}{
					"target": NewL0LoadBalancerLogic(mockLogic.Logic()),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0LoadBalancerLogic)
				logic.DeleteLoadBalancer("lbid")

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(1, len(tags))
				reporter.AssertEqual(tags[0].EntityID, "not_lbid")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.DeleteLoadBalancer error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteLoadBalancer(gomock.Any()).
					Return(fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if err := logic.DeleteLoadBalancer(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteLoadBalancer(gomock.Any()).
					Return(nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if err := logic.DeleteLoadBalancer(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateLoadBalancer(t *testing.T) {
	request := models.CreateLoadBalancerRequest{
		EnvironmentID:    "envid",
		LoadBalancerName: "lb_name",
		IsPublic:         true,
		Ports: []models.Port{
			models.Port{
				HostPort:      80,
				ContainerPort: 80,
				Protocol:      "tcp",
			},
		},
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should error if request.EnvironmentID is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				request := models.CreateLoadBalancerRequest{LoadBalancerName: "lb_name"}
				if _, err := logic.CreateLoadBalancer(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should error if request.LoadBalancerName is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				request := models.CreateLoadBalancerRequest{EnvironmentID: "envid"}
				if _, err := logic.CreateLoadBalancer(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should error if name/environment_id tags already exist",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "lbid",
					EntityType: "load_balancer",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "lbid",
					EntityType: "load_balancer",
					Key:        "environment_id",
					Value:      "envid",
				})

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				request := models.CreateLoadBalancerRequest{
					EnvironmentID:    "envid",
					LoadBalancerName: "some_name",
				}

				if _, err := logic.CreateLoadBalancer(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should call backend.CreateLoadBalancer with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					CreateLoadBalancer("lb_name", "envid", true, request.Ports).
					Return(&models.LoadBalancer{}, nil)

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)
				logic.CreateLoadBalancer(request)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.CreateLoadBalancer error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					CreateLoadBalancer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if _, err := logic.CreateLoadBalancer(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should add correct name tag in database",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				loadBalancer := &models.LoadBalancer{
					EnvironmentID:    "envid",
					LoadBalancerID:   "lbid",
					LoadBalancerName: "lb_name",
				}

				mockLogic.Backend.EXPECT().
					CreateLoadBalancer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(loadBalancer, nil)

				mockLogic.UseSQLite(t)

				return map[string]interface{}{
					"target": NewL0LoadBalancerLogic(mockLogic.Logic()),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				// todo: setup id generation

				logic := testMap["target"].(*L0LoadBalancerLogic)
				if _, err := logic.CreateLoadBalancer(request); err != nil {
					reporter.Error(err)
				}

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(len(tags), 2)

				nameTag := models.EntityTag{
					EntityID:   "lbid",
					EntityType: "load_balancer",
					Key:        "name",
					Value:      "lb_name",
				}

				environmentTag := models.EntityTag{
					EntityID:   "lbid",
					EntityType: "load_balancer",
					Key:        "environment_id",
					Value:      "envid",
				}

				reporter.AssertInSlice(nameTag, tags)
				reporter.AssertInSlice(environmentTag, tags)
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if _, err := logic.CreateLoadBalancer(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestUpdateLoadBalancer(t *testing.T) {
	request := models.UpdateLoadBalancerRequest{
		Ports: []models.Port{
			models.Port{
				HostPort:      80,
				ContainerPort: 80,
				Protocol:      "tcp",
			},
		},
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.UpdateLoadBalancer with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					UpdateLoadBalancer("lbid", request.Ports).
					Return(&models.LoadBalancer{}, nil)

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)
				logic.UpdateLoadBalancer("lbid", request.Ports)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.UpdateLoadBalancer error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					UpdateLoadBalancer(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if _, err := logic.UpdateLoadBalancer("", request.Ports); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					UpdateLoadBalancer(gomock.Any(), gomock.Any()).
					Return(&models.LoadBalancer{}, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0LoadBalancerLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0LoadBalancerLogic)

				if _, err := logic.UpdateLoadBalancer("", request.Ports); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
