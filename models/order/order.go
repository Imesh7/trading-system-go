package order

import (
	"fmt"
	"trading-system-go/database"
	"github.com/gofiber/fiber/v2"
)

type Order struct {
	OrderId     int         `gorm:"AUTO_INCREMENT;PRIMARY_KEY;not null" json:"order_id"`
	UserId      int32       `json:"user_id"`
	OrderType   OrderType   `json:"order_type"`
	Type        string      `json:"type"`
	Price       float64     `json:"price"`
	Volume      float64     `json:"volume"`
	BuyingPair  string      `json:"buying_pair"`
	SellingPair string      `json:"selling_pair"`
	OrderStatus OrderStatus `json:"order_status"`
	CreatedAt   int64       `gorm:"autoCreateTime" json:"created_at"`
}

type OrderType int

const (
	_ OrderType = iota
	MarketOrderBuy
	LimitOrderBuy
	MarketOrderSell
	LimitOrderSell
)

type OrderStatus string

const (
	PartialyFilled OrderStatus = "PartialyFilled"
	Filled         OrderStatus = "Filled"
	NotFilled      OrderStatus = "UnFilled"
)


//create a order
func (order *Order) CreateOrder(c *fiber.Ctx) error {

	err := c.BodyParser(&order)
	if err != nil {
		return c.Status(400).JSON("json error")
	}

	if order.UserId == 0 || (order.OrderType == 0 || order.OrderType > 4) || order.BuyingPair == "" || order.SellingPair == "" || order.Price == 0 || order.Volume == 0 {
		s := fmt.Sprintf("Unprocessable Entity %d", order.OrderType)
		return c.Status(402).JSON(s)
	}

	//first create the order in DB
	order.OrderStatus = "UnFilled"
	database.DB.DataBase.Create(&order)
	if err != nil {
		return c.Status(400).JSON("something went wrong on database cannot create a order")
	}

	OrderMatchProducer("topic",order.OrderId)

	return c.Status(200).JSON(&order)
}



//get all the order details
func (order *Order) GetOrders(c *fiber.Ctx) error {
	var orderList []Order

	database.DB.DataBase.Find(&orderList)
	if database.DB.DataBase.Error != nil {
		fmt.Println(database.DB.DataBase.Error)
		return c.Status(400).JSON("error")
	}
	return c.Status(200).JSON(&orderList)
}


