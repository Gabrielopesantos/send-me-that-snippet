package paste

import (
	"fmt"
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/model"
	utls "github.com/gabrielopesantos/smts/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
)

type pasteHandlers struct {
	dbConn *gorm.DB
	cfg    *config.Config
}

func NewHandlers(dbConn *gorm.DB, cfg *config.Config) *pasteHandlers {
	return &pasteHandlers{
		cfg:    cfg,
		dbConn: dbConn,
	}
}

func (h *pasteHandlers) Filter() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		filter := model.Paste{}
		err := ctx.QueryParser(&filter)
		log.Printf("%+v", filter)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		var results []model.Paste
		// Breaks this functions
		err = h.dbConn.Where(&filter, "expired").Find(&results).Error
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		err = ctx.JSON(results)
		if err != nil {
			log.Printf("%v", err)
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
			ctx.Response().SetBodyString("Failed to parse request body")
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		err = val.Struct(&paste)
		if err != nil {
			ctx.Response().SetBodyString("Invalid request body")
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		paste.Id = utls.RandSeq(12)

		err = h.dbConn.Create(&paste).Error
		if err != nil {
			ctx.Response().SetBodyString("Failed to register paste")
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.JSON(paste)
	}
}
func (h *pasteHandlers) Get() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pId := ctx.Params("pId", "")

		paste := model.Paste{}
		err := h.dbConn.First(&paste, "id  = ?", pId).Error
		if err != nil {
			ctx.Response().SetBodyString("Not found")
			return ctx.SendStatus(fiber.StatusNotFound)
		}

		return ctx.JSON(paste)
	}
}

func (h *pasteHandlers) Delete() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pId := ctx.Params("pId")

		paste := model.Paste{}
		err := h.dbConn.Delete(&paste, "id = ?", pId)
		if err != nil {
			ctx.Response().SetBodyString("Not found")
			return ctx.SendStatus(fiber.StatusNotFound)
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
			ctx.Response().SetBodyString("Nothing to update")
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		err = h.dbConn.Where("id = ?", pId).Updates(paste).Error
		if err != nil {
			ctx.Response().SetBodyString("Failed to update item")
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
