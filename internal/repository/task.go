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
		On("CONFLICT (external_id) DO UPDATE").
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

func (r *TaskRepository) MarkTasksAsProcessed(ctx context.Context, taskIDs []int64) error {
	if len(taskIDs) == 0 {
		return nil
	}

	_, err := r.db.NewUpdate().
		Model(&model.Task{}).
		Set("processed = ?", true).
		Where("id IN (?)", bun.In(taskIDs)).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("error marking tasks as processed: %v", err)
	}

	return nil
}
