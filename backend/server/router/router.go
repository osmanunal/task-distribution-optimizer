package router

import (
	"task-distribution-optimizer/internal/repository"
	"task-distribution-optimizer/internal/service"
	"task-distribution-optimizer/server/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func Setup(app *fiber.App, db *bun.DB) {
	// Repository'leri oluştur
	taskRepo := repository.NewTaskRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)

	// Service'leri oluştur
	taskService := service.NewTaskService(nil, taskRepo, employeeRepo)

	// Handler'ları oluştur
	taskHandler := handler.NewTaskHandler(taskService)

	// API gruplarını oluştur
	api := app.Group("/api")

	// Task endpoint'leri
	tasks := api.Group("/tasks")
	tasks.Get("/plan", taskHandler.TaskPlanner)
}
