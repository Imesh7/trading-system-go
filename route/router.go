package router

import (
	"trading-system-go/models/order"
	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {
	app.Post("/create-order", func(c *fiber.Ctx) error {
		var order order.Order
		return order.CreateOrder(c)
	})

	app.Get("/get-orders", func(c *fiber.Ctx) error {
		var order order.Order
		return order.GetOrders(c)
	})
}
