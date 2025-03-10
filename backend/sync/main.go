package main

import (
	"context"
	"fmt"
	"os"
	"task-distribution-optimizer/internal/model"
	"task-distribution-optimizer/internal/provider"
	"task-distribution-optimizer/internal/repository"
	"task-distribution-optimizer/internal/service"
	"task-distribution-optimizer/pkg/config"
	"task-distribution-optimizer/pkg/database"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "sync",
		Commands: commands,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

var commands = []*cli.Command{
	{
		Name:  "start",
		Usage: "task senkronizasyonunu başlat",
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			cfg := config.Read()
			db := database.ConnectDB(cfg.DBConfig)

			taskRepo := repository.NewTaskRepository(db)
			taskProvider := provider.NewProvider1()
			syncService := service.NewTaskService(taskProvider, taskRepo, nil)

			fmt.Println("Task senkronizasyonu başlatılıyor...")
			if err := syncService.SyncTasks(ctx); err != nil {
				fmt.Printf("Senkronizasyon hatası: %v", err)
				return err
			}

			fmt.Println("Task senkronizasyon başarılı!")
			return nil
		},
	},
	{
		Name:  "plan",
		Usage: "görevleri çalışanlara dağıt",
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			cfg := config.Read()
			db := database.ConnectDB(cfg.DBConfig)

			taskRepo := repository.NewTaskRepository(db)
			employeeRepo := repository.NewEmployeeRepository(db)
			plannerService := service.NewTaskService(nil, taskRepo, employeeRepo)

			result, err := plannerService.TaskPlanner(ctx)
			if err != nil {
				fmt.Printf("Görev dağıtım hatası: %v", err)
				return err
			}

			fmt.Println("En hızlı bitecek atama planı:")
			fmt.Printf("Toplam süre: %d hafta\n", result.TotalWeeks)
			fmt.Println("==========================================")
			fmt.Println("Haftalık çalışan programı:")

			for _, workload := range result.Workloads {
				fmt.Printf("Çalışan: %s (ID: %d)\n", workload.EmployeeName, workload.EmployeeID)
				fmt.Printf("  Toplam iş yükü: %d saat\n", workload.TotalHours)

				if workload.TotalHours == 0 {
					fmt.Printf("  Atanan görevler: Yok\n")
				} else {
					fmt.Printf("  Atanan görevler:\n")
					for _, assignment := range workload.Assignments {
						fmt.Printf("    - %s (ID: %d): %d saat\n",
							assignment.TaskName, assignment.TaskExternalID, assignment.Duration)
					}

					fmt.Printf("  Haftalık plan:\n")
					for _, week := range workload.WeeklyPlan {
						fmt.Printf("    Hafta %d: %d saat çalışma\n", week.WeekNumber, week.Hours)
					}
				}
				fmt.Println("------------------------------------------")
			}

			fmt.Println("Görev dağıtımı başarıyla tamamlandı!")
			return nil
		},
	},
	{
		Name:  "add-emp",
		Usage: "çalışanları ekle",
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			cfg := config.Read()
			db := database.ConnectDB(cfg.DBConfig)

			employeeRepo := repository.NewEmployeeRepository(db)

			employees := []model.Employee{
				{Name: "DEV1", Difficulty: 1},
				{Name: "DEV2", Difficulty: 2},
				{Name: "DEV3", Difficulty: 3},
				{Name: "DEV4", Difficulty: 4},
				{Name: "DEV5", Difficulty: 5},
			}

			for _, emp := range employees {
				if err := employeeRepo.CreateEmployee(ctx, emp); err != nil {
					fmt.Printf("Çalışan eklenirken hata oluştu %s: %v\n", emp.Name, err)
					return err
				}
				fmt.Printf("Çalışan eklendi: %s (Zorluk: %dx)\n", emp.Name, emp.Difficulty)
			}

			fmt.Println("Tüm çalışanlar başarıyla eklendi!")
			return nil
		},
	},
}
