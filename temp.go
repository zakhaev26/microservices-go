package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func mainx() {
	var wg sync.WaitGroup
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1:9092",
		"acks":              "all",
		// "client.id":         "something",
	})
	cons_ch := make(chan string)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	delivery_chan := make(chan kafka.Event, 10000)
	topic := "Topic"

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9092",
			"group.id":          "cons-1",
			"auto.offset.reset": "smallest",
		})
		cons_ch <- "Cons bangaya"
		if err != nil {
			log.Fatal(err)
		}
		err = consumer.Subscribe(topic, nil)
		if err != nil {
			log.Fatal(err)
		}

		for {
			ev := consumer.Poll(100)

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Consumed message : %+s\n", string(e.Value))
			case *kafka.Error:
				fmt.Printf("%v\n", e)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var i int = 0
		for {
			msg := "Hi Soubhik! " + strconv.Itoa(i)
			i++
			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(msg)},
				delivery_chan,
			)

			if err != nil {
				fmt.Printf("Failed to produce message: %v\n", err)
				os.Exit(1)
			}
			e := <-delivery_chan
			fmt.Println("Delivered", e.String())
			time.Sleep(time.Second * 3)
		}
	}()

	msg := <-cons_ch
	fmt.Println(msg)

	wg.Wait()
}
