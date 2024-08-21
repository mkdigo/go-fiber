package main

import (
	"encoding/json"
	"net/http"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	app := fiber.New()

	app.Get("/", public)

	app.Get("/joke", func(c *fiber.Ctx) error {
		var apiResponse map[string]string

		response, err := http.Get("https://api.chucknorris.io/jokes/random")

		if err != nil {
			c.Status(fiber.StatusInternalServerError).SendString("request chucknorris api error")
		}

		json.NewDecoder(response.Body).Decode(&apiResponse)

		return c.SendString(apiResponse["value"])
	})

	app.Post("/login", login)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte("super_secret"),
		},
	}))

	app.Get("/restricted", restricted)

	app.Listen(":3000")
}

func public(c *fiber.Ctx) error {
	return c.SendString("This is public route")
}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("super_secret"))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome" + name)
}
