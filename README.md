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

To run the application, use the following command:
make sure you already installed docker


```bash
docker-compose up
```

## See Orderbook Data
To visulalize orderbook data go to redis commander:
If you already run the `docker-compose up` will works, browser below url

```bash
http://localhost:8081/
```
