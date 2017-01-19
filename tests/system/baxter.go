package system

import (
	"testing"
)

type Baxter struct {
	URL string
}

func NewBaxter(t *testing.T, url string) *Baxter {
	return &Baxter{
		URL: url,
	}
}

func (b *Baxter) Die() {

}
