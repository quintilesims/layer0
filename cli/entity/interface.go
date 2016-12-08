package entity

import (
	"github.com/quintilesims/layer0/cli/printer/table"
)

type Entity interface {
	Table() table.Table
}
