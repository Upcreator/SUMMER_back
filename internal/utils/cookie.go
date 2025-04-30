package utils

import (
	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/gofiber/fiber/v2"
	"time"
)

func NewCookie(name string, value string, path string, expires time.Time) *fiber.Cookie {
	if path == "" {
		path = "/"
	}

	secure := false

	if initializers.AppConfig.Stand == "prod" {
		secure = true
	}

	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HTTPOnly: true,
		Secure:   secure,
		Domain:   initializers.AppConfig.CookieDomain,
		SameSite: "Lax",
	}
}
