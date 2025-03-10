package viewmodel

import (
	"task-distribution-optimizer/internal/port"
)

// WeeklyWorkViewModel haftalık çalışma planını temsil eder
type WeeklyWorkViewModel struct {
	WeekNumber int `json:"week_number"`
	Hours      int `json:"hours"`
}

// AssignmentViewModel, bir çalışana atanan görevleri temsil eder
type AssignmentViewModel struct {
	TaskID         int64  `json:"task_id"`
	TaskName       string `json:"task_name"`
	TaskExternalID int64  `json:"task_external_id"`
	Duration       int    `json:"duration"` // Saat cinsinden
}

// EmployeeWorkloadViewModel, bir çalışanın iş yükünü temsil eder
type EmployeeWorkloadViewModel struct {
	EmployeeID   int64                 `json:"employee_id"`
	EmployeeName string                `json:"employee_name"`
	Difficulty   int                   `json:"difficulty"`
	TotalHours   int                   `json:"total_hours"`
	Assignments  []AssignmentViewModel `json:"assignments"`
	WeeklyPlan   []WeeklyWorkViewModel `json:"weekly_plan"`
}

// TaskDistributionResponse, API yanıtında kullanılacak veri modelini temsil eder
type TaskDistributionResponse struct {
	TotalWeeks int                         `json:"total_weeks"`
	Workloads  []EmployeeWorkloadViewModel `json:"workloads"`
}

// ToViewModel, iç modeli API yanıt modeline dönüştürür
func (vm TaskDistributionResponse) ToViewModel(m port.TaskDistributionResult) TaskDistributionResponse {
	response := TaskDistributionResponse{
		TotalWeeks: m.TotalWeeks,
		Workloads:  make([]EmployeeWorkloadViewModel, 0, len(m.Workloads)),
	}

	for _, workload := range m.Workloads {
		employeeVM := EmployeeWorkloadViewModel{
			EmployeeID:   workload.EmployeeID,
			EmployeeName: workload.EmployeeName,
			Difficulty:   workload.Difficulty,
			TotalHours:   workload.TotalHours,
			Assignments:  make([]AssignmentViewModel, 0, len(workload.Assignments)),
			WeeklyPlan:   make([]WeeklyWorkViewModel, 0, len(workload.WeeklyPlan)),
		}

		// Görevleri dönüştür
		for _, assignment := range workload.Assignments {
			employeeVM.Assignments = append(employeeVM.Assignments, AssignmentViewModel{
				TaskID:         assignment.TaskID,
				TaskName:       assignment.TaskName,
				TaskExternalID: assignment.TaskExternalID,
				Duration:       assignment.Duration,
			})
		}

		// Haftalık planı dönüştür
		for _, weekly := range workload.WeeklyPlan {
			employeeVM.WeeklyPlan = append(employeeVM.WeeklyPlan, WeeklyWorkViewModel{
				WeekNumber: weekly.WeekNumber,
				Hours:      weekly.Hours,
			})
		}

		response.Workloads = append(response.Workloads, employeeVM)
	}

	return response
}
