package repository

import (
	"context"
	"fmt"
	"task-distribution-optimizer/internal/model"
	"task-distribution-optimizer/internal/port"

	"github.com/uptrace/bun"
)

type TaskRepository struct {
	db bun.IDB
}

func NewTaskRepository(db *bun.DB) port.TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) UpsertTasks(ctx context.Context, tasks []model.Task) error {
	if len(tasks) == 0 {
		return nil
	}

	_, err := r.db.NewInsert().
		Model(&tasks).
		On("CONFLICT (name, external_id) DO UPDATE").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("tasks upsert error: %v", err)
	}

	return nil
}

func (r *TaskRepository) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task

	err := r.db.NewSelect().
		Model(&tasks).
		Order("id ASC").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error fetching all tasks: %v", err)
	}

	return tasks, nil
}
