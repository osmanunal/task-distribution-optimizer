package repository

import (
	"context"
	"fmt"
	"task-distribution-optimizer/internal/model"
	"task-distribution-optimizer/internal/port"

	"github.com/uptrace/bun"
)

type EmployeeRepository struct {
	db bun.IDB
}

func NewEmployeeRepository(db *bun.DB) port.EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) CreateEmployee(ctx context.Context, employee model.Employee) error {
	_, err := r.db.NewInsert().Model(&employee).Exec(ctx)
	return err
}

func (r *EmployeeRepository) GetAllEmployees(ctx context.Context) ([]model.Employee, error) {
	var employees []model.Employee

	err := r.db.NewSelect().
		Model(&employees).
		Order("id ASC").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error fetching all employees: %v", err)
	}

	return employees, nil
}
