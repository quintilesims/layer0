package main

import (
	"context"

	"github.com/quintilesims/layer0/cli/client"
)

type Layer0Client struct {
	API         client.Client
	StopContext context.Context
}
