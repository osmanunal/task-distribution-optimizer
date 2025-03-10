package service

import (
	"context"
	"log"
	"math"
	"sort"
	"sync"
	"task-distribution-optimizer/internal/model"
	"task-distribution-optimizer/internal/port"
	"task-distribution-optimizer/pkg/utils"
)

type TaskService struct {
	taskProvider port.TaskProvider
	taskRepo     port.TaskRepository
	employeeRepo port.EmployeeRepository
}

func NewTaskService(taskProvider port.TaskProvider, taskRepo port.TaskRepository, employeeRepo port.EmployeeRepository) port.TaskService {
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
func (s *TaskService) TaskPlanner(ctx context.Context) (port.TaskDistributionResult, error) {
	const weeklyWorkHours = 45 // Haftalık çalışma saati limiti

	tasks, err := s.taskRepo.GetAllTasks(ctx)
	if err != nil {
		return port.TaskDistributionResult{}, err
	}

	employees, err := s.employeeRepo.GetAllEmployees(ctx)
	if err != nil {
		return port.TaskDistributionResult{}, err
	}

	if len(tasks) == 0 || len(employees) == 0 {
		return port.TaskDistributionResult{}, err
	}

	// Çalışanları ve görevleri zorluk derecesine göre sırala
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].Difficulty > employees[j].Difficulty
	})

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Difficulty > tasks[j].Difficulty
	})

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

	// Görevleri en uygun çalışanlara dağıt
	for _, task := range tasks {
		bestEmployeeIndex := -1
		minTotalHours := math.MaxInt32
		minTaskDuration := math.MaxInt32

		for i, emp := range employees {
			duration := emp.ComputeTaskDuration(task)

			// Çalışan iş yükü ve görev süresi kombinasyonuna göre en iyi eşleşmeyi bul
			// İş yükü düşük olanları ve görevi hızlı yapabilenleri tercih et
			if workloads[i].TotalHours < minTotalHours ||
				(workloads[i].TotalHours == minTotalHours && duration < minTaskDuration) {
				minTotalHours = workloads[i].TotalHours
				minTaskDuration = duration
				bestEmployeeIndex = i
			}
		}

		if bestEmployeeIndex != -1 {
			emp := employees[bestEmployeeIndex]
			duration := emp.ComputeTaskDuration(task)

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
			workloads[bestEmployeeIndex].TotalHours += duration
		}
	}

	// Haftalık planları paralel olarak hesapla
	var wg sync.WaitGroup
	maxWeeksCh := make(chan int, len(workloads))

	for i := range workloads {
		if workloads[i].TotalHours > 0 {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				workloads[idx].WeeklyPlan = calculateWeeklyPlan(workloads[idx].TotalHours, weeklyWorkHours)

				weeksNeeded := (workloads[idx].TotalHours + weeklyWorkHours - 1) / weeklyWorkHours // ceiling division
				maxWeeksCh <- weeksNeeded
			}(i)
		}
	}

	// Tüm go routine'ler tamamlandığında channel'ı kapat
	go func() {
		wg.Wait()
		close(maxWeeksCh)
	}()

	maxWeeks := 0
	for weeks := range maxWeeksCh {
		if weeks > maxWeeks {
			maxWeeks = weeks
		}
	}

	result := port.TaskDistributionResult{
		TotalWeeks: maxWeeks,
		Workloads:  workloads,
	}

	return result, nil
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
