package paste

import (
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/model"
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

func (h *pasteHandlers) Insert() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var err error
		val := validator.New()
		paste := model.Paste{}
		err = ctx.BodyParser(&paste)
		if err != nil {
			log.Print(err)
			return err
		}

		err = val.Struct(&paste)
		if err != nil {
			log.Print(err)
			return err
		}

		h.dbConn.Create(&paste)

		return ctx.JSON(paste)
	}
}
func (h *pasteHandlers) Get() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pId := ctx.Params("pId")

		paste := model.Paste{}
		h.dbConn.First(&paste, "id  = ?", pId)

		return ctx.JSON(paste)
	}
}

func (h *pasteHandlers) Delete() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		pId := ctx.Params("pId")

		paste := model.Paste{}
		h.dbConn.Delete(&paste, "id = ?", pId)

		return ctx.JSON(paste)
	}
}
