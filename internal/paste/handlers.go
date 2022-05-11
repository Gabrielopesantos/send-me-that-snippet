package paste

import (
	"fmt"
	"github.com/gabrielopesantos/smts/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
		fmt.Print("Insert")
		return nil
	}
}
func (h *pasteHandlers) Get() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fmt.Print("Get")
		return nil
	}
}

func (h *pasteHandlers) Delete() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fmt.Print("Delete")
		return nil
	}
}
