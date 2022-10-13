package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"time"
)

type BaseResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

func PerformAPICall(c *fiber.Ctx) error {
	client := resty.New()
	client.SetDebug(true)
	client.SetContentLength(true)
	client.SetTimeout(time.Duration(time.Duration(60) * time.Second))

	response, err := client.R().
		Get("https://random-data-api.com/api/v2/users")

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			BaseResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error",
				Data: &fiber.Map{
					"content": err.Error(),
				},
			},
		)
	}

	return c.Status(http.StatusCreated).JSON(
		BaseResponse{
			Status:  http.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"content": gjson.Parse(response.String()).Value(),
			},
		},
	)
}

func main() {
	app := fiber.New()

	app.Get("/random-data", func(c *fiber.Ctx) error {
		return PerformAPICall(c)
	})

	app.Listen(":6000")
}
