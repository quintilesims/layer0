package logic

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/commmon/db"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestGetCertificate(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.GetCertificate with correct param",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					GetCertificate("crt_id").
					Return(&models.Certificate{}, nil)

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)
				logic.GetCertificate("crt_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.GetCertificate error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					GetCertificate(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				if _, err := logic.GetCertificate("crt_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate model with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetCertificate(gomock.Any()).
					Return(&models.Certificate{CertificateID: "crt_id"}, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "crt_id",
					EntityType: "certificate",
					Key:        "name",
					Value:      "some_name",
				})

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				certificate, err := logic.GetCertificate("crt_id")
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(certificate.CertificateID, "crt_id")
				reporter.AssertEqual(certificate.CertificateName, "some_name")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetCertificate(gomock.Any()).
					Return(&models.Certificate{}, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				if _, err := logic.GetCertificate("crt_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListCertificates(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.ListCertificates",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListCertificates().
					Return([]*models.Certificate{}, nil)

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)
				logic.ListCertificates()
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.ListCertificates error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListCertificates().
					Return(nil, fmt.Errorf("some error"))

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				if _, err := logic.ListCertificates(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate models with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				certificates := []*models.Certificate{
					&models.Certificate{CertificateID: "crt_id_1"},
					&models.Certificate{CertificateID: "crt_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListCertificates().
					Return(certificates, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "crt_id_1",
					EntityType: "certificate",
					Key:        "name",
					Value:      "some_name_1",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "crt_id_2",
					EntityType: "certificate",
					Key:        "name",
					Value:      "some_name_2",
				})

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				certificates, err := logic.ListCertificates()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(certificates[0].CertificateID, "crt_id_1")
				reporter.AssertEqual(certificates[0].CertificateName, "some_name_1")
				reporter.AssertEqual(certificates[1].CertificateID, "crt_id_2")
				reporter.AssertEqual(certificates[1].CertificateName, "some_name_2")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				certificates := []*models.Certificate{
					&models.Certificate{CertificateID: "crt_id_1"},
					&models.Certificate{CertificateID: "crt_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListCertificates().
					Return(certificates, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				if _, err := logic.ListCertificates(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteCertificate(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.DeleteCertificate with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					DeleteCertificate("crt_id").
					Return(nil)

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)
				logic.DeleteCertificate("crt_id")
			},
		},
		testutils.TestCase{
			Name: "Should delete certificate tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteCertificate(gomock.Any()).
					Return(nil)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "crt_id",
					EntityType: "certificate",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "not_crt_id",
					EntityType: "certificate",
					Key:        "name",
					Value:      "some_name",
				})

				return map[string]interface{}{
					"target": NewL0CertificateLogic(mockLogic.Logic()),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0CertificateLogic)
				logic.DeleteCertificate("crt_id")

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(1, len(tags))
				reporter.AssertEqual(tags[0].EntityID, "not_crt_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.DeleteCertificate error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteCertificate(gomock.Any()).
					Return(fmt.Errorf("some error"))

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				if err := logic.DeleteCertificate(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteCertificate(gomock.Any()).
					Return(nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				if err := logic.DeleteCertificate(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateCertificate(t *testing.T) {
	request := models.CreateCertificateRequest{
		CertificateName:  "some_name",
		PublicCert:       "public",
		PrivateKey:       "private",
		IntermediateCert: "intermed",
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should error if request.CertificateName is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				if _, err := logic.CreateCertificate(models.CreateCertificateRequest{}); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should call backend.CreateCertificate with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					CreateCertificate("some_name", "public", "private", "intermed").
					Return(&models.Certificate{}, nil)

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)
				logic.CreateCertificate(request)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.CreateCertificate error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					CreateCertificate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				if _, err := logic.CreateCertificate(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should add correct name tag in database",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					CreateCertificate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Certificate{CertificateID: "crt_id"}, nil)

				mockLogic.UseSQLite(t)

				return map[string]interface{}{
					"target": NewL0CertificateLogic(mockLogic.Logic()),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0CertificateLogic)
				if _, err := logic.CreateCertificate(request); err != nil {
					reporter.Error(err)
				}

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(len(tags), 1)
				reporter.AssertEqual(tags[0].EntityID, "crt_id")
				reporter.AssertEqual(tags[0].EntityType, "certificate")
				reporter.AssertEqual(tags[0].Key, "name")
				reporter.AssertEqual(tags[0].Value, "some_name")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					CreateCertificate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Certificate{}, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0CertificateLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0CertificateLogic)

				if _, err := logic.CreateCertificate(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
