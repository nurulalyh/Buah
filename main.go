package main

import (
	"fmt"
	"net/http"
)

type Buah struct {
	Id int `json: "id"`
	Name string `json: "name"`
	Price float64 `json: "price"`
}
type Response struct {
	Code int `json: "code"`
	Message string `json: "message"`
	Data interface{} `json: "data"`
}

func main() {
	http.HandleFunc("/api/v1/buah-sehat", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})

	port := "8000"
	fmt.Println("Server Run On Port", port)
	http.ListenAndServe(":"+port, nil)
}