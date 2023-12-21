package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	db "github.com/zakhaev26/microservices-go/database"
	model "github.com/zakhaev26/microservices-go/models"
)

func HandleDataPost(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var incomingScore model.AdminData

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := json.NewDecoder(r.Body).Decode(&incomingScore)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()

		ins, err := db.Collection.InsertOne(context.TODO(), incomingScore)
		fmt.Println("Pubbed :", ins)
		if err != nil {
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(ins)
		}
	}()

	delivery_chan := make(chan kafka.Event, 10000)
	topic := "GC"

	wg.Add(1)
	go func() {
		defer wg.Done()

		p, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": "127.0.0.1:9092",
			"acks":              "all",
		})

		if err != nil {
			fmt.Printf("Failed to produce message: %v\n", err)
		}

		jsonPayload, err := json.Marshal(incomingScore)
		if err != nil {
			fmt.Println(err)
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(string(jsonPayload))},
			delivery_chan,
		)

		if err != nil {
			fmt.Printf("Failed to produce message: %v\n", err)
			os.Exit(1)
		}
		<-delivery_chan
	}()

	wg.Wait()
}
