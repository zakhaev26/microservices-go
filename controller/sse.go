package controller

import (
	"fmt"
	"net/http"
	"time"
)

func SSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("REQ")
	eventChannel := make(chan string)

	go func() {
		for {
			eventChannel <- "Channel se ara "
			time.Sleep(time.Second * 3)
		}
	}()

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
