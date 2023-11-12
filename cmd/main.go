package main

import (
	"fmt"
	"log"
	"trading-system-go/database"
	"trading-system-go/models/balance"
	"trading-system-go/models/order"
	router "trading-system-go/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Print("aplication started.........................")
	app := fiber.New()
	database.ConnectDatabase()
	database.ConnectToRedis()
	database.DB.DataBase.AutoMigrate(&order.Order{}, &balance.Balance{})
	router.AppRoutes(app)
	go order.OrderMatchConsumer("topic")
	log.Fatal(app.Listen(":8000"))
}
