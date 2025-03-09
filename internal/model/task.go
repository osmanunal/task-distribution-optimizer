package model

import (
	"taskmanager/pkg/model"
)

type Task struct {
	model.BaseModel
	ExternalID int64  `bun:",notnull"`
	Name       string `bun:",notnull"`
	Difficulty int    `bun:",notnull"`
	Duration   int    `bun:",notnull"`
}

type TaskProviderResponse struct {
	ExternalID int64
	Difficulty int
	Duration   int
	Name       string
}
