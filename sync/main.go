package main

import (
	"context"
	"log"
	"os"
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
			taskProvider := provider.NewProvider1("https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-one")
			syncService := service.NewTaskSyncService(taskProvider, taskRepo, nil)

			log.Println("Task senkronizasyonu başlatılıyor...")
			if err := syncService.SyncTasks(ctx); err != nil {
				log.Printf("Senkronizasyon hatası: %v", err)
				return err
			}

			log.Println("Task senkronizasyon başarılı!")
			return nil
		},
	},
}
