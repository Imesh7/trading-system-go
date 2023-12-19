package routes

import (
	"fmt"
	"trading-system-go/internal/data/order"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {
	app.Post("/create-order", CreateOrderRoute)
	app.Get("/get-orders", GetOrdersRoute)
}

func CreateOrderRoute(c *fiber.Ctx) error {
	var orderData order.Order
	err := c.BodyParser(&orderData)
	if err != nil {
		return c.Status(402).JSON(&orderData)
	}
	fmt.Println()
	if orderData.UserId == 0 || (orderData.OrderType <= 0 || orderData.OrderType > 4) || orderData.BuyingPair == "" || orderData.SellingPair == "" || orderData.Price == 0 || orderData.Volume == 0 {
		s := fmt.Sprintf("Unprocessable Entity %d", orderData.OrderType)
		return c.Status(402).JSON(s)
	}
	order, err := order.CreateOrder(&orderData)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON("something went wrong on database cannot create a order")
	}
	return c.Status(200).JSON(&order)
}

func GetOrdersRoute(c *fiber.Ctx) error {
	orderList, err := order.GetOrders(c)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	return c.Status(200).JSON(orderList)
}
