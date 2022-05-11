package middleware

import (
	"crypto/subtle"
	"encoding/base64"

	"github.com/gabrielopesantos/smts/config"
	"github.com/gofiber/fiber/v2"
)

type MiddlewareManager struct {
	cfg          *config.Config
	expectedAuth string
}

func NewMiddlewareManager(cfg *config.Config) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:          cfg,
		expectedAuth: "Basic" + " " + base64.StdEncoding.EncodeToString([]byte(cfg.ServerConfig.BasicAuthUser+":"+cfg.ServerConfig.BasicAuthPassword)),
	}
}

func (m *MiddlewareManager) BasicAuthMiddleware(next fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		providedAuth := ctx.Get("Authorization")

		authMatch := subtle.ConstantTimeCompare([]byte(m.expectedAuth), []byte(providedAuth)) == 1
		if authMatch {
			return next(ctx)
		}

		ctx.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		ctx.Status(fiber.StatusUnauthorized)
		return ctx.SendString("Unauthorized")
	}
}

//func (m *MiddlewareManager) BasicAuthMiddlewareUse() fiber.Handler {
//	return func(ctx *fiber.Ctx) error {
//		providedAuth := ctx.Get("Authorization")
//
//		authMatch := subtle.ConstantTimeCompare([]byte(m.expectedAuth), []byte(providedAuth)) == 1
//		if authMatch {
//			return ctx.Next()
//		}
//
//		ctx.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
//		ctx.Status(fiber.StatusUnauthorized)
//		return ctx.SendString("Unauthorized")
//	}
//}
