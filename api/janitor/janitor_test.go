package janitor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJanitorRun(t *testing.T) {
	var called bool
	janitor := NewJanitor("", func() error {
		called = true
		return nil
	})

	janitor.Run()
	assert.True(t, called)
}

func TestJanitorRunEvery(t *testing.T) {
	c := make(chan bool)
	janitor := NewJanitor("", func() error {
		c <- true
		return nil
	})

	ticker := janitor.RunEvery(time.Nanosecond)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
		<-c
	}
}
