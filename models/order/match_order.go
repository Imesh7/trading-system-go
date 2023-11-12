package order

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"trading-system-go/database"
)

func MatchOrder(orderId string) {
	var order Order
	///fmt.Println("Match order got the value from kafka consumer", orderId)
	//first get the unfilled or partfilled order with opposite
	//order type of this particular order
	database.DB.DataBase.Where("order_id = ?", orderId).First(&order)
	//try to match with the price.
	///fmt.Fprintln(os.Stdout, []any{"order details order_id=%d , order_type=%s order_price=%v", order.OrderId, order.Type, order.Price}...)
	//if the price match , remove from the orderbook & update the order status filled
	ctx := context.Background()
	buyOrderkey := fmt.Sprintf("buy-orders")
	sellOrderkey := fmt.Sprintf("sell-orders")
	//Buy order
	if order.Type == "bid" {
		//first try to match with the redis after that get that order & validate it &
		//remove from the order & create a sell order against,
		sellOrderList := GetOrderBook(ctx, sellOrderkey)

		if len(sellOrderList) == 0 {
			//add a new buy order to the orderbook
			AddToOrderBook(ctx, buyOrderkey, order.OrderId, serializeOrder(order))
		} else {
			//sort the order with the price low to high
			sort.SliceStable(sellOrderList, func(i, j int) bool {
				//if price is equal sort with order createdAt(first fill the oldest order)
				if sellOrderList[i].Price == sellOrderList[j].Price {
					return sellOrderList[i].CreatedAt < sellOrderList[j].CreatedAt
				}
				return sellOrderList[i].Price < sellOrderList[j].Price
			})

			for _, v := range sellOrderList {
				//if bid price is equal or higher than the sell order in the order book
				if order.Price >= v.Price {
					if order.Volume <= v.Volume {
						//set the buy order status as fully filled
						//update the postgres order as filled

						//set the redis & postgres sell order to or if equal fully filled / partially filled
						balanceVolume := v.Volume - order.Volume
						if balanceVolume == 0 {
							//update the orderbook in fully filled
							v.OrderStatus = Filled
							/* err := UpdateOrderBook(ctx, sellOrderkey, v.OrderId, serializeOrder(v))
							if err != nil {
								fmt.Println(err)
								} */
							//remove item from the order book
							errs := RemoveFromOrderBook(ctx, sellOrderkey, v.OrderId)
							if errs != nil {
								fmt.Println(errs)
							}
						} else {
							//update the orderbook in partially filled
							v.OrderStatus = PartialyFilled
							v.Volume = v.Volume - order.Volume
							UpdateOrderBook(ctx, sellOrderkey, v.OrderId, serializeOrder(v))
						}
						break
					} else {
						//buy order should be partially filled
						//sell fully filled
						v.OrderStatus = Filled
						//update the orderbook in partially filled
						/* err := UpdateOrderBook(ctx, sellOrderkey, v.OrderId, serializeOrder(v))
						if err != nil {
							fmt.Println(err)
							} */
						order.Volume = order.Volume - v.Volume
						errs := RemoveFromOrderBook(ctx, sellOrderkey, v.OrderId)
						if errs != nil {
							fmt.Println(errs)
						}
					}
				} else {
					//add a new buy order to the orderbook
					AddToOrderBook(ctx, buyOrderkey, order.OrderId, serializeOrder(order))
					break
				}
			}
		}
	} else if order.Type == "ask" {
		//first try to match with the redis after that get that order & validate it &
		//remove from the order & create a sell order against,
		buyOrderList := GetOrderBook(ctx, buyOrderkey)

		if len(buyOrderList) == 0 {
			AddToOrderBook(ctx, sellOrderkey, order.OrderId, serializeOrder(order))
		} else {
			//sort the order with the price low to high
			sort.SliceStable(buyOrderList, func(i, j int) bool {
				//if price is equal sort with order createdAt(first fill the oldest order)
				if buyOrderList[i].Price == buyOrderList[j].Price {
					return buyOrderList[i].CreatedAt < buyOrderList[j].CreatedAt
				}
				return buyOrderList[i].Price > buyOrderList[j].Price
			})

			for _, v := range buyOrderList {
				//if bid price is equal or higher than the sell order in the order book
				if order.Price <= v.Price {
					if order.Volume <= v.Volume {
						//set the buy order status as fully filled
						//update the postgres order as filled

						//set the redis & postgres sell order to or if equal fully filled / partially filled
						balanceVolume := v.Volume - order.Volume
						if balanceVolume == 0 {
							//update the orderbook in fully filled
							v.OrderStatus = Filled
							/* err := UpdateOrderBook(ctx, buyOrderkey, v.OrderId, serializeOrder(v))
							if err != nil {
								fmt.Println(err)
							} */
							//remove item from the order book
							errs := RemoveFromOrderBook(ctx, buyOrderkey, v.OrderId)
							if errs != nil {
								fmt.Println(errs)
							}
						} else {
							//update the orderbook in partially filled
							v.OrderStatus = PartialyFilled
							v.Volume = v.Volume - order.Volume
							UpdateOrderBook(ctx, buyOrderkey, v.OrderId, serializeOrder(v))
						}
						break
					} else {
						//buy order should be partially filled
						//sell fully filled
						v.OrderStatus = Filled
						//update the orderbook in partially filled
						/* err := UpdateOrderBook(ctx, buyOrderkey, v.OrderId, serializeOrder(v))
						if err != nil {
							fmt.Println(err)
						} */
						order.Volume = order.Volume - v.Volume
						errs := RemoveFromOrderBook(ctx, buyOrderkey, v.OrderId)
						if errs != nil {
							fmt.Println(errs)
						}

						//recurrsion
					}
				} else {
					//add a new buy order to the orderbook
					AddToOrderBook(ctx, sellOrderkey, order.OrderId, serializeOrder(order))
					break
				}
			}
		}
	}
}

func serializeOrder(order Order) string {
	data, _ := json.Marshal(order)
	return string(data)
}
