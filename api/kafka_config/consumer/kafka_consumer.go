package order_kafka_consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"trading-system-go/api/match_order"

	"github.com/IBM/sarama"
)

func OrderMatchConsumer(topic string) {
	config := sarama.NewConfig()
	config.ClientID = "go-kafka-consumer"
	kafkaHost := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	consumer, err := sarama.NewConsumer([]string{kafkaHost}, config)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"Errors is %s", err}...)
		log.Fatal(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			valueString := string(msg.Value)
			fmt.Fprintln(os.Stdout, []any{"Received integer value: %s", valueString}...)
			match_order.MatchOrder(valueString)
		case <-signals:
			return
		}
	}
}
