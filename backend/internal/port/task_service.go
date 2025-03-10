package port

import (
	"context"
)

// Bir görevin bir çalışana atanmasını temsil eder
type Assignment struct {
	TaskID         int64
	TaskName       string
	TaskExternalID int64
	Duration       int
}

// Bir çalışanın iş programını temsil eder
type EmployeeWorkload struct {
	EmployeeID   int64
	EmployeeName string
	Difficulty   int
	TotalHours   int
	Assignments  []Assignment
	WeeklyPlan   []WeeklyWork
}

// Haftalık çalışma planını temsil eder
type WeeklyWork struct {
	WeekNumber int
	Hours      int
}

// Çıktı modelini temsil eder
type TaskDistributionResult struct {
	TotalWeeks int
	Workloads  []EmployeeWorkload
}

// TaskService, task senkronizasyon servisi için port
type TaskService interface {
	// SyncTasks, görevleri senkronize eder
	SyncTasks(ctx context.Context) error

	// TaskPlanner, görevleri çalışanlara planlar ve en optimal dağılımı hesaplar
	TaskPlanner(ctx context.Context) (TaskDistributionResult, error)
}
