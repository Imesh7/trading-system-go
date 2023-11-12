# Trading system using go

Create a Order process in a trading system.
 - Create a order
 - Maintain orderbook
 - Match the order with orderbook

 Tech stack used
 --------------------
 - Backend server
 - Postgres (db)
 - Kafka
 - Redis (orderbook)

## Running the Application

To run the application, use the following command.  
Make sure you already installed Docker.


```bash
docker-compose up
```

## See Orderbook Data
To visulalize orderbook data go to redis commander:
If you already run the `docker-compose up` will works, browser below url

```bash
http://localhost:8081/
```

```json

http://localhost:8000/create-order
//create a Ask(sell) order
{
    "user_id":1000000,
    "order_type":4,
    "type":"ask",
    "price":100,
    "volume":150,
    "buying_pair":"usd",
    "selling_pair":"btc"
}

//create a Bid(buy) order
{
    "user_id":1000000,
    "order_type":4,
    "type":"bid",
    "price":100,
    "volume":150,
    "buying_pair":"usd",
    "selling_pair":"btc"
}

//Get orders 
http://localhost:8000/get-orders
```

