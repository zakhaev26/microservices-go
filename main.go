package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zakhaev26/microservices-go/controller"
	db "github.com/zakhaev26/microservices-go/database"
)

func main() {
	db.Init()
	r := mux.NewRouter()
	r.HandleFunc("/post", controller.HandleDataPost)
	fmt.Println("Running at Port 8080")
	http.ListenAndServe(":8080", r)
}
