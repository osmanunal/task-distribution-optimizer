package service

import (
	"context"
	"log"
	"math"
	"sort"
	"task-distribution-optimizer/internal/model"
	"task-distribution-optimizer/internal/port"
	"task-distribution-optimizer/pkg/utils"
)

type TaskService struct {
	taskProvider port.TaskProvider
	taskRepo     port.TaskRepository
	employeeRepo port.EmployeeRepository
}

func NewTaskService(taskProvider port.TaskProvider, taskRepo port.TaskRepository, employeeRepo port.EmployeeRepository) port.TaskSyncService {
	return &TaskService{
		taskProvider: taskProvider,
		taskRepo:     taskRepo,
		employeeRepo: employeeRepo,
	}
}

// SyncTasks, task'ları senkronize eder
func (s *TaskService) SyncTasks(ctx context.Context) error {
	tasks, err := s.taskProvider.GetTasks(ctx)
	if err != nil {
		return err
	}

	var modelTasks []model.Task
	for _, t := range tasks {
		modelTasks = append(modelTasks, model.Task{
			ExternalID: t.ExternalID,
			Difficulty: t.Difficulty,
			Duration:   t.Duration,
			Name:       t.Name,
		})
	}

	err = s.taskRepo.UpsertTasks(ctx, modelTasks)
	if err != nil {
		return err
	}

	log.Printf("Task senkronizasyonu tamamlandı")
	return nil
}

// TaskPlanner, görevleri çalışanlara en uygun şekilde dağıtır ve sonucu döndürür
func (s *TaskService) TaskPlanner(ctx context.Context) (*port.TaskDistributionResult, error) {
	const weeklyWorkHours = 45 // Haftalık çalışma saati limiti

	// Görev ve çalışanları getir
	tasks, err := s.taskRepo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	employees, err := s.employeeRepo.GetAllEmployees(ctx)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 || len(employees) == 0 {
		return &port.TaskDistributionResult{
			TotalWeeks: 0,
			Workloads:  []port.EmployeeWorkload{},
		}, nil
	}

	// Görevleri zorluklarına göre sırala
	sort.Slice(tasks, func(i, j int) bool {
		// Her görev için ortalama süreyi hesapla
		var totalDurationI, totalDurationJ int
		for _, emp := range employees {
			totalDurationI += emp.ComputeTaskDuration(tasks[i])
			totalDurationJ += emp.ComputeTaskDuration(tasks[j])
		}
		avgDurationI := float64(totalDurationI) / float64(len(employees))
		avgDurationJ := float64(totalDurationJ) / float64(len(employees))
		return avgDurationI > avgDurationJ // Zorları önce ata
	})

	// Çalışan iş yüklerini hazırla
	workloads := make([]port.EmployeeWorkload, len(employees))
	for i, emp := range employees {
		workloads[i] = port.EmployeeWorkload{
			EmployeeID:   emp.ID,
			EmployeeName: emp.Name,
			Difficulty:   emp.Difficulty,
			TotalHours:   0,
			Assignments:  []port.Assignment{},
			WeeklyPlan:   []port.WeeklyWork{},
		}
	}

	// Her çalışanın toplam süresini takip et
	completionTimes := make([]int, len(employees))

	// Her görevi en uygun çalışana ata
	for _, task := range tasks {
		bestEmployeeIndex := -1
		earliestCompletionTime := math.MaxInt32

		// Bu görevi en erken tamamlayacak çalışanı bul
		for empIndex, emp := range employees {
			duration := emp.ComputeTaskDuration(task)
			completionTime := completionTimes[empIndex] + duration

			if completionTime < earliestCompletionTime {
				earliestCompletionTime = completionTime
				bestEmployeeIndex = empIndex
			}
		}

		// En uygun çalışana görevi ata
		if bestEmployeeIndex != -1 {
			duration := employees[bestEmployeeIndex].ComputeTaskDuration(task)

			// Çalışanın toplam süresini güncelle
			completionTimes[bestEmployeeIndex] += duration
			workloads[bestEmployeeIndex].TotalHours += duration

			// Atamayı kaydet
			assignment := port.Assignment{
				TaskID:         task.ID,
				TaskName:       task.Name,
				TaskExternalID: task.ExternalID,
				Duration:       duration,
			}
			workloads[bestEmployeeIndex].Assignments = append(
				workloads[bestEmployeeIndex].Assignments,
				assignment,
			)
		}
	}

	// Her çalışanın haftalık planını hesapla
	for i := range workloads {
		workloads[i].WeeklyPlan = calculateWeeklyPlan(workloads[i].TotalHours, weeklyWorkHours)
	}

	// Toplam hafta sayısını hesapla
	totalWeeks := s.calculateTotalWeeks(workloads, weeklyWorkHours)

	return &port.TaskDistributionResult{
		TotalWeeks: totalWeeks,
		Workloads:  workloads,
	}, nil
}

// calculateWeeklyPlan, toplam saati haftalık plana dönüştürür
func calculateWeeklyPlan(totalHours, weeklyWorkHours int) []port.WeeklyWork {
	var plan []port.WeeklyWork

	remainingHours := totalHours
	weekNumber := 1

	for remainingHours > 0 {
		weeklyHours := utils.Min(remainingHours, weeklyWorkHours)

		plan = append(plan, port.WeeklyWork{
			WeekNumber: weekNumber,
			Hours:      weeklyHours,
		})

		remainingHours -= weeklyHours
		weekNumber++
	}

	return plan
}

// calculateTotalWeeks, en uzun süren çalışana göre toplam hafta sayısını hesaplar
func (s *TaskService) calculateTotalWeeks(workloads []port.EmployeeWorkload, weeklyWorkHours int) int {
	maxWeeks := 0

	for _, workload := range workloads {
		weeks := (workload.TotalHours + weeklyWorkHours - 1) / weeklyWorkHours // Ceiling division
		if weeks > maxWeeks {
			maxWeeks = weeks
		}
	}

	return maxWeeks
}
