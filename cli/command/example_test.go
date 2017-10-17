package command

import (
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/config"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	args := Args{"one", "two"}
	flags := Flags{
		"name":  "dev",
		"count": 1,
	}

	c := NewContext(t, args, flags)

	assert.Equal(t, "one", c.Args().Get(0))
	assert.Equal(t, "two", c.Args().Get(1))
	assert.Equal(t, 1, c.Int("count"))
	assert.Equal(t, "dev", c.String("name"))

	c = NewContext(t, nil, nil,
		SetVersion("latest"),
		SetNoWait(true),
		SetTimeout(time.Minute*3))

	assert.Equal(t, "latest", c.App.Version)
	assert.Equal(t, true, c.GlobalBool(config.FLAG_NO_WAIT))
	assert.Equal(t, time.Minute*3, c.GlobalDuration(config.FLAG_TIMEOUT))
}
