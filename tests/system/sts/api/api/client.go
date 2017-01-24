package sts

import (
	"fmt"
	"github.com/dghubble/sling"
	"testing"
)

type SystemTestService struct {
	T   *testing.T
	URL string
}

func NewSystemTestService(t *testing.T, url string) *SystemTestService {
	return &SystemTestService{
		T:   t,
		URL: fmt.Sprintf("http://%s", url),
	}
}

func (s *SystemTestService) sling() *sling.Sling {
	return sling.New().Base(s.URL)
}

func (s *SystemTestService) Die() {
	resp, err := s.sling().Post("/health").BodyJSON("").ReceiveSuccess("")
	if err != nil {
		s.T.Fatal(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		s.T.Fatalf("STS returned invalid status code :%s", resp.Status)
	}
}
