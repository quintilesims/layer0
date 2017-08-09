package main

import (
	"context"

	"github.com/quintilesims/layer0/common/client"
)

type Layer0Client struct {
	API         client.Client
	StopContext context.Context
}
