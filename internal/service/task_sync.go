package service

import (
	"context"
	"log"
	"sort"
	"task-distribution-optimizer/internal/model"

	"task-distribution-optimizer/internal/port"
)

type TaskService struct {
	taskProvider port.TaskProvider
	taskRepo     port.TaskRepository
	employeeRepo port.EmployeeRepository
}

func NewTaskSyncService(taskProvider port.TaskProvider, taskRepo port.TaskRepository, employeeRepo port.EmployeeRepository) port.TaskSyncService {
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

func (s *TaskService) TaskPlanner(ctx context.Context) error {
	tasks, err := s.taskRepo.GetAllTasks(ctx)
	if err != nil {
		return err
	}

	employees, err := s.employeeRepo.GetAllEmployees(ctx)
	if err != nil {
		return err
	}

	if len(tasks) == 0 || len(employees) == 0 {
		log.Printf("Görev planlaması için yeterli veri yok: %d görev, %d çalışan", len(tasks), len(employees))
		return nil
	}

	// Haftalık çalışma saati limiti
	const weeklyWorkHours = 45

	// Görevleri zorluk derecesine göre sırala (en zor olanlar önce)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Difficulty > tasks[j].Difficulty
	})

	// Her çalışanın görev listesi ve toplam çalışma saati
	type EmployeeSchedule struct {
		Employee       model.Employee
		Tasks          []model.Task
		TotalHours     int
		WeeklySchedule map[int][]model.Task // Hafta numarası -> o haftanın görevleri
		WeeklyHours    map[int]int          // Hafta numarası -> o haftanın toplam saati
	}

	// Çalışanları zorluk derecelerine göre sırala (en yetenekliler önce)
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].Difficulty > employees[j].Difficulty
	})

	schedules := make([]EmployeeSchedule, len(employees))
	for i, emp := range employees {
		schedules[i] = EmployeeSchedule{
			Employee:       emp,
			Tasks:          []model.Task{},
			WeeklySchedule: make(map[int][]model.Task),
			WeeklyHours:    make(map[int]int),
		}
	}

	// Her görevi, en uygun çalışana atama
	var updatedTasks []model.Task
	for _, task := range tasks {
		// En az yüklü olan çalışanı bul
		bestEmployeeIdx := 0
		minTotalHours := computeTaskDuration(task, schedules[0].Employee) + schedules[0].TotalHours

		for i := 1; i < len(schedules); i++ {
			// Bu çalışanın görev üzerinde çalışacağı süre
			workHours := computeTaskDuration(task, schedules[i].Employee)
			totalHours := workHours + schedules[i].TotalHours

			if totalHours < minTotalHours {
				minTotalHours = totalHours
				bestEmployeeIdx = i
			}
		}

		// Görevi en uygun çalışana ata
		taskCopy := task
		empID := schedules[bestEmployeeIdx].Employee.ID
		taskCopy.EmployeeID = &empID

		schedules[bestEmployeeIdx].Tasks = append(schedules[bestEmployeeIdx].Tasks, taskCopy)
		schedules[bestEmployeeIdx].TotalHours += computeTaskDuration(task, schedules[bestEmployeeIdx].Employee)

		updatedTasks = append(updatedTasks, taskCopy)
	}

	// Haftalık programı oluştur
	maxWeeks := 0
	for i := range schedules {
		week := 0
		weeklyHours := 0

		for _, task := range schedules[i].Tasks {
			// Bu görevin tamamlanma süresi
			taskHours := computeTaskDuration(task, schedules[i].Employee)

			// Eğer bu görev mevcut haftaya sığmıyorsa, sonraki haftaya geç
			if weeklyHours+taskHours > weeklyWorkHours {
				week++
				weeklyHours = 0
			}

			// Görevi bu haftaya ekle
			if schedules[i].WeeklySchedule[week] == nil {
				schedules[i].WeeklySchedule[week] = []model.Task{}
			}
			schedules[i].WeeklySchedule[week] = append(schedules[i].WeeklySchedule[week], task)

			// Haftalık saatleri güncelle
			schedules[i].WeeklyHours[week] = schedules[i].WeeklyHours[week] + taskHours
			weeklyHours += taskHours
		}

		// Toplam hafta sayısını güncelle
		totalWeeks := week
		if weeklyHours > 0 {
			totalWeeks++
		}

		if totalWeeks > maxWeeks {
			maxWeeks = totalWeeks
		}
	}

	// Veritabanını güncelle
	err = s.taskRepo.UpsertTasks(ctx, updatedTasks)
	if err != nil {
		return err
	}

	// Sonuçları logla
	log.Printf("\n--- GÖREV PLANLAMA SONUÇLARI ---")
	log.Printf("İşlerin minimum tamamlanma süresi: %d hafta", maxWeeks)

	for _, schedule := range schedules {
		log.Printf("\nÇalışan: %s (Zorluk: %dx)", schedule.Employee.Name, schedule.Employee.Difficulty)
		log.Printf("Toplam atanan çalışma saati: %d saat", schedule.TotalHours)

		for week := 0; week < maxWeeks; week++ {
			tasks := schedule.WeeklySchedule[week]
			hours := schedule.WeeklyHours[week]

			if len(tasks) > 0 {
				log.Printf("  Hafta %d (%d saat):", week+1, hours)
				for _, task := range tasks {
					taskHours := computeTaskDuration(task, schedule.Employee)
					log.Printf("    - %s (Süre: %d saat, Zorluk: %d)", task.Name, taskHours, task.Difficulty)
				}
			} else {
				log.Printf("  Hafta %d: Görev yok", week+1)
			}
		}
	}

	return nil
}

// computeTaskDuration, bir görevin belirli bir çalışan tarafından tamamlanma süresini hesaplar
func computeTaskDuration(task model.Task, employee model.Employee) int {
	// Çalışanın zorluk derecesi 0 ise, varsayılan olarak 1 kullan
	difficulty := employee.Difficulty
	if difficulty <= 0 {
		difficulty = 1
	}

	// Formül: Görev süresi / Çalışan zorluk derecesi
	// Örnek: 10 saatlik 1 zorluk işi, 2x zorluktaki developer 5 saatte bitirebilir
	return task.Duration / difficulty
}
