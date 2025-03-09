package port

import (
	"context"
	"task-distribution-optimizer/internal/model"
)

// TaskRepository, task veritabanı işlemleri için port
type TaskRepository interface {
	// UpsertTasks, task'ları günceller veya yeni ekler
	UpsertTasks(ctx context.Context, tasks []model.Task) error

	// GetAllTasks, tüm task'ları getirir
	GetAllTasks(ctx context.Context) ([]model.Task, error)

	// MarkTasksAsProcessed, belirtilen task'ları işlenmiş olarak işaretler
	MarkTasksAsProcessed(ctx context.Context, taskIDs []int64) error
}
