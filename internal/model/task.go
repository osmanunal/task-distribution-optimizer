package model

import (
	"task-distribution-optimizer/pkg/model"
)

type Task struct {
	model.BaseModel
	ExternalID int64  `bun:",notnull"`
	Name       string `bun:",notnull"`
	Difficulty int    `bun:",notnull"`
	Duration   int    `bun:",notnull"`
	EmployeeID *int64
	Employee   Employee `bun:"rel:belongs-to,join:employee_id=id"`
}

type TaskProviderResponse struct {
	ExternalID int64
	Difficulty int
	Duration   int
	Name       string
}
