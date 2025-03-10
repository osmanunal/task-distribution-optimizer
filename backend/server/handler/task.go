package handler

import (
	"task-distribution-optimizer/internal/port"
	"task-distribution-optimizer/server/viewmodel"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	taskService port.TaskService
}

func NewTaskHandler(taskService port.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) TaskPlanner(c *fiber.Ctx) error {
	taskDistributionResult, err := h.taskService.TaskPlanner(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	vm := viewmodel.TaskDistributionResponse{}
	response := vm.ToViewModel(taskDistributionResult)
	return c.Status(fiber.StatusOK).JSON(response)
}
