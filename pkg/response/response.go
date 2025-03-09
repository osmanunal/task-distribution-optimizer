package response

import (
	"task-distribution-optimizer/pkg/errorx"
	"task-distribution-optimizer/pkg/model"

	"github.com/gofiber/fiber/v2"
)

func Success(ctx *fiber.Ctx, data interface{}, dataCount ...int) error {
	count := 0
	if len(dataCount) > 0 {
		count = dataCount[0]
	}
	response := model.Response{
		Data:      data,
		DataCount: count,
		Error:     "",
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func Error(ctx *fiber.Ctx, errType int, err ...string) error {
	var err1 string
	if errType == fiber.StatusInternalServerError {
		err1 = errorx.InternalServerError
		if len(err) > 0 {
			err1 = err[0]
		}
	} else if errType == fiber.StatusNotFound {
		err1 = errorx.NotFound
		if len(err) > 0 {
			err1 = err[0]
		}
	} else if errType == fiber.StatusBadRequest {
		err1 = errorx.BadRequest
		if len(err) > 0 {
			err1 = err[0]
		}
	}

	return ctx.Status(errType).JSON(model.Response{
		Data:      nil,
		DataCount: 0,
		Error:     err1,
	})
}
