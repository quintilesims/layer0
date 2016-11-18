package entity

import (
	"gitlab.imshealth.com/xfra/layer0/cli/printer/table"
)

type Entity interface {
	Table() table.Table
}
