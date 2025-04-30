package middleware

import (
	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"slices"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return []byte(initializers.AppConfig.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid token claims",
		})
	}

	c.Locals("userId", claims["sub"])
	c.Locals("userRole", claims["role"])
	return c.Next()
}

func RoleRequired(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем роль пользователя из локальных данных контекста
		userRole := c.Locals("userRole").(string)

		// Проверяем, содержится ли роль пользователя в списке необходимых ролей
		if !slices.Contains(requiredRoles, userRole) {
			// Если роль пользователя не входит в список, возвращаем ошибку 403 (Forbidden)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You are not allowed to access this resource",
			})
		}

		// Если роль подходит, продолжаем выполнение следующего middleware или обработчика
		return c.Next()
	}
}
