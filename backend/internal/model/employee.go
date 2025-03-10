package model

import "task-distribution-optimizer/pkg/model"

type Employee struct {
	model.BaseModel

	Name       string `bun:",notnull"`
	Difficulty int    `bun:",notnull"`
}

func (e *Employee) ComputeTaskDuration(task Task) int {
	return task.Difficulty * task.Duration / e.Difficulty
}
