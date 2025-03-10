package port

import (
	"context"
	"task-distribution-optimizer/internal/model"
)

type TaskProvider interface {
	GetTasks(ctx context.Context) ([]model.TaskProviderResponse, error)
}
