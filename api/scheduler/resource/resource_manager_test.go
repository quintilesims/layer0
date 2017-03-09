package resource

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// test cases:
// not enough ports, need to scale up
// not enough memory, need to scale up
// not enough memory per-instance, but would work if you did a bad calculation, should scale up

// we have way too much room
// we have exactly enough room

// we have enough room if you
// we have barely too little room
// we have way too little room
func TestResourceManager_shouldScaleUp(t *testing.T) {

	assert.Equal(t, 0, 0)
}
