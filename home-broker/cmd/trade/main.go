package main

import (
	"encoding/json"
	"fmt"
	"home-broker/internal/infra/kafka"
	"home-broker/internal/market/dto"
	"home-broker/internal/market/entity"
	"home-broker/internal/market/transformer"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func Main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	}

	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	go kafka.Consume(kafkaMsgChan) // T2

	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() //T3

	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Print(string(msg.Value))

			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)

			if err != nil {
				panic(err)
			}

			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", "  ")

		if err != nil {
			fmt.Print(err)
		}

		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
