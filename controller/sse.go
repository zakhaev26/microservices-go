package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("REQ")
	eventChannel := make(chan string)

	consumer1, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "cons-1",
		"auto.offset.reset": "smallest",
	})

	if err != nil {
		log.Fatal(err)
	}

	err = consumer1.Subscribe("GC", nil)

	if err != nil {
		log.Fatal(err)
	}

	//Generation Part of Code
	go func() {
		for {

			ev := consumer1.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Consumed message : %+s\n", string(e.Value))
				eventChannel <- string(e.Value)
			case *kafka.Error:
				fmt.Printf("%v\n", e)
			}

			// eventChannel <- "Channel se ara "
			// time.Sleep(time.Second * 3)
		}
	}()

	//SSE Sending Part Of code
	for {
		select {
		case msg, ok := <-eventChannel:
			if !ok {
				fmt.Println("Channel Closed")
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush()
			fmt.Println("Sent")
		case <-r.Context().Done():
			fmt.Println("Client disconnected")
			return
		}
	}

}
