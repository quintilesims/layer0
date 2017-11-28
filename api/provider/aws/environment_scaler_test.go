package aws

import (
	"testing"
)

func TestResourceManagerScaleUp_noProviders(t *testing.T) {
	// there are 0 providers in the cluster
	// there is 1 consumer
	// we should scale up to size 1

}
