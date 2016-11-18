package ecsbackend

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs/id"
	"gitlab.imshealth.com/xfra/layer0/common/aws/iam"
	"gitlab.imshealth.com/xfra/layer0/common/aws/iam/mock_iam"
	"gitlab.imshealth.com/xfra/layer0/common/testutils"
	"testing"
)

type MockECSCertificateManager struct {
	IAM *mock_iam.MockProvider
}

func NewMockECSCertificateManager(ctrl *gomock.Controller) *MockECSCertificateManager {
	return &MockECSCertificateManager{
		IAM: mock_iam.NewMockProvider(ctrl),
	}
}

func (this *MockECSCertificateManager) Certificate() *ECSCertificateManager {
	return NewECSCertificateManager(this.IAM)
}

func TestGetCertificate(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call iam.ListCertificates with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				mockCertificate.IAM.EXPECT().
					ListCertificates(CertificatePath()).
					Return(nil, nil)

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				certificate := target.(*ECSCertificateManager)
				certificate.GetCertificate("crtid")
			},
		},
		testutils.TestCase{
			Name: "Should return layer0-formatted certificate id",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				ecsCertificateID := id.L0CertificateID("crtid").ECSCertificateID()
				certificate := iam.NewServerCertificateMetadata(ecsCertificateID.String(), "")

				mockCertificate.IAM.EXPECT().
					ListCertificates(gomock.Any()).
					Return([]*iam.ServerCertificateMetadata{certificate}, nil)

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSCertificateManager)

				certificate, err := manager.GetCertificate("crtid")
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual("crtid", certificate.CertificateID)
			},
		},
		testutils.TestCase{
			Name: "Should propagate iam.ListCertificates error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				mockCertificate.IAM.EXPECT().
					ListCertificates(gomock.Any()).
					Return(nil, fmt.Errorf("some_error"))

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSCertificateManager)

				if _, err := manager.GetCertificate("crtid"); err == nil {
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
			Name: "Should call iam.ListCertificates with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				mockCertificate.IAM.EXPECT().
					ListCertificates(CertificatePath()).
					Return(nil, nil)

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				certificate := target.(*ECSCertificateManager)
				certificate.ListCertificates()
			},
		},
		testutils.TestCase{
			Name: "Should return layer0-formatted certificate ids",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				ecsCertificateID := id.L0CertificateID("crtid").ECSCertificateID()
				metadata := iam.NewServerCertificateMetadata(ecsCertificateID.String(), "")
				mockCertificate.IAM.EXPECT().
					ListCertificates(gomock.Any()).
					Return([]*iam.ServerCertificateMetadata{metadata}, nil)

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSCertificateManager)

				certificates, err := manager.ListCertificates()
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(len(certificates), 1)
				reporter.AssertEqual(certificates[0].CertificateID, "crtid")
			},
		},
		testutils.TestCase{
			Name: "Should propagate iam.ListCertificates error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				mockCertificate.IAM.EXPECT().
					ListCertificates(gomock.Any()).
					Return(nil, fmt.Errorf("some_error"))

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSCertificateManager)

				if _, err := manager.ListCertificates(); err == nil {
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
			Name: "Should call iam.DeleteServerCertificates with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				ecsCertificateID := id.L0CertificateID("crtid").ECSCertificateID()
				mockCertificate.IAM.EXPECT().
					DeleteServerCertificate(ecsCertificateID.String()).
					Return(nil)

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				certificate := target.(*ECSCertificateManager)
				certificate.DeleteCertificate("crtid")
			},
		},
		testutils.TestCase{
			Name: "Should propagate iam.DeleteServerCertificate error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				mockCertificate.IAM.EXPECT().
					DeleteServerCertificate(gomock.Any()).
					Return(fmt.Errorf("some_error"))

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSCertificateManager)

				if err := manager.DeleteCertificate("crtid"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateCertificate(t *testing.T) {
	defer id.StubIDGeneration("crtid")()

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call iam.UploadServerCertificates with proper params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				ecsCertificateID := id.L0CertificateID("crtid").ECSCertificateID()
				metadata := iam.NewServerCertificateMetadata(ecsCertificateID.String(), "")

				mockCertificate.IAM.EXPECT().
					UploadServerCertificate(
						ecsCertificateID.String(),
						CertificatePath(),
						"public",
						"private",
						stringp("chain"),
					).Return(metadata, nil)

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				certificate := target.(*ECSCertificateManager)
				certificate.CreateCertificate("some_name", "public", "private", "chain")
			},
		},
		testutils.TestCase{
			Name: "Should propagate iam.UploadServerCertificate error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				mockCertificate.IAM.EXPECT().
					UploadServerCertificate(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).Return(nil, fmt.Errorf("some_error"))

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSCertificateManager)

				if _, err := manager.CreateCertificate("some_name", "public", "private", "chain"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should return layer0-formatted certificate id",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCertificate := NewMockECSCertificateManager(ctrl)

				ecsCertificateID := id.L0CertificateID("crtid").ECSCertificateID()
				metadata := iam.NewServerCertificateMetadata(ecsCertificateID.String(), "")

				mockCertificate.IAM.EXPECT().
					UploadServerCertificate(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any()).
					Return(metadata, nil)

				return mockCertificate.Certificate()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSCertificateManager)

				certificate, err := manager.CreateCertificate("some_name", "public", "private", "chain")
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(certificate.CertificateID, "crtid")
			},
		},
	}

	testutils.RunTests(t, testCases)
}
