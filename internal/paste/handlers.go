package paste

import (
	"fmt"
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/model"
	"github.com/gabrielopesantos/smts/pkg/logger"
	utls "github.com/gabrielopesantos/smts/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type pasteHandlers struct {
	dbConn Repository
	logger *logger.Logger
	cfg    *config.Config
}

func NewHandlers(dbConn Repository, logger *logger.Logger, cfg *config.Config) *pasteHandlers {
	return &pasteHandlers{
		cfg:    cfg,
		logger: logger,
		dbConn: dbConn,
	}
}

func (h *pasteHandlers) Filter() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		filter := model.Paste{}
		err := ctx.QueryParser(&filter)

		if err != nil {
			h.logger.Error(fmt.Sprintf("Filter - %s", err.Error()), nil)
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		results, err := h.dbConn.Filter(&filter)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Filter - %s", err.Error()), nil)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		err = ctx.JSON(results)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Filter - %s", err.Error()), nil)
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (h *pasteHandlers) Insert() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var err error
		val := validator.New()
		paste := model.Paste{}
		err = ctx.BodyParser(&paste)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Insert - %s", err.Error()), nil)
			return ctx.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
		}

		err = val.Struct(&paste)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Insert - %s", err.Error()), nil)
			return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request boyd")
		}

		paste.Id = utls.RandSeq(12)

		err = h.dbConn.Insert(&paste)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Insert - %s", err.Error()), nil)
			return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to register paste")
		}

		return ctx.JSON(paste)
	}
}
func (h *pasteHandlers) Get() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pId := ctx.Params("pId", "")

		paste, err := h.dbConn.Get(pId)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Get - %s", err.Error()), nil)
			return ctx.Status(fiber.StatusNotFound).SendString("Not found")
		}

		return ctx.JSON(paste)
	}
}

func (h *pasteHandlers) Delete() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pId := ctx.Params("pId")

		err, paste := h.dbConn.Delete(pId)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Get - %v", err), nil)
			return ctx.Status(fiber.StatusNotFound).SendString("Not found")
		}

		return ctx.JSON(paste)
	}
}

func (h *pasteHandlers) Update() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pId := ctx.Params("pId")

		paste := model.Paste{}
		err := ctx.BodyParser(&paste)
		fmt.Printf("Something - %s | %+v\n", pId, paste)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Update - %s", err.Error()), nil)
			return ctx.Status(fiber.StatusBadRequest).SendString("Nothing to update")
		}

		err = h.dbConn.Update(pId, &paste)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Update - %s", err.Error()), nil)
			return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to update item")
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
