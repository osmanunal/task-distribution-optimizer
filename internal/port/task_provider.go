package port

import (
	"context"
	"taskmanager/internal/model"
)

type TaskProvider interface {
	GetTasks(ctx context.Context) ([]model.TaskProviderResponse, error)
}
