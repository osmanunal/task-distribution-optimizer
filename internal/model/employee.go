package model

import "task-distribution-optimizer/pkg/model"

type Employee struct {
	model.BaseModel

	Name       string `bun:",notnull"`
	Difficulty int    `bun:",notnull"`
	Workload   int    `bun:",notnull"`
}
