package port

import (
	"context"
	"task-distribution-optimizer/internal/model"
)

// TaskRepository, task veritabanı işlemleri için port
type EmployeeRepository interface {
	GetAllEmployees(ctx context.Context) ([]model.Employee, error)
	CreateEmployee(ctx context.Context, employee model.Employee) error
}
