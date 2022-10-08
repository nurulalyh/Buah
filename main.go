package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	models "github.com/nurulalyh/Buah/models"
	"github.com/nurulalyh/Buah/repository"
)

func sendResponse(code int, message string, data interface{}, w http.ResponseWriter) {
	resp := models.Response{
		Code:    code,
		Data:    data,
		Message: message,
	}
	dataByte, err := json.Marshal(resp)

	if err != nil {
		resp := models.Response{
			Code:    http.StatusInternalServerError,
			Data:    nil,
			Message: "Internal Server Error",
		}
		dataByte, _ = json.Marshal(resp)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dataByte)
}

var db *sql.DB

// func GetFruit(w http.ResponseWriter, r *http.Request) {

// }

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:rootroot@localhost/fruits?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}

	if err = db.Ping(); err != nil {
		panic(err.Error())
	}
	defer db.Close()

	r := mux.NewRouter()

	//Method Get
	r.HandleFunc("/api/v1/fruits", func(w http.ResponseWriter, r *http.Request) {
		fruits, err := repository.GetFruits(db)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "Internal Server Error, Get Fruits", err.Error(), w)
			return
		}

		sendResponse(http.StatusOK, "Success", fruits, w)
	}).Methods(http.MethodGet)

	//Method Put
	r.HandleFunc("/api/v1/fruits/{id}", func(w http.ResponseWriter, r *http.Request) {
		//get query param
		id := mux.Vars(r)["id"]

		if id == "" {
			sendResponse(http.StatusBadRequest, "Bad Request, Data Id Params Is Null", nil, w)
			return
		}

		// idInt, err := strconv.Atoi(id)
		// rows, err := db.Query("select id, name, price from fruits where id = $1", id)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "Internal Server Error, Get Fruits", nil, w)
			return
		}

		fruit, err := repository.GetFruit(db, id)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "Internal Server Error, Get Fruits return error", err.Error(), w)
			return
		}

		if fruit.Id == 0 {
			if err != nil {
				sendResponse(http.StatusBadRequest, "bad request", nil, w)
			}
		}
		dataByte, err := io.ReadAll(r.Body)

		if err != nil {
			sendResponse(http.StatusBadRequest, "bad request", nil, w)
			return
		}
		defer r.Body.Close()

		var newFruit models.Fruit
		err = json.Unmarshal(dataByte, &newFruit)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), w)
			return
		}

		fruit.Name = newFruit.Name
		fruit.Price = newFruit.Price

		err = repository.UpdateFruits(db, fruit)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), w)
			return
		}

		sendResponse(http.StatusOK, "Success Update", nil, w)
	}).Methods(http.MethodPut)

	//Method Post
	r.HandleFunc("/api/v1/fruits", func(w http.ResponseWriter, r *http.Request) {
		dataByte, err := io.ReadAll(r.Body)
		if err != nil {
			sendResponse(http.StatusBadRequest, "Bad Request", nil, w)
		}
		defer r.Body.Close()

		var fruit models.Fruit
		err = json.Unmarshal(dataByte, &fruit)

		if err != nil {
			sendResponse(http.StatusInternalServerError, "Internal Server Error", nil, w)
		}

		err = repository.CreateFruits(db, fruit)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), w)
		}

		sendResponse(http.StatusCreated, "Success", nil, w)
		return
	}).Methods(http.MethodPost)

	//Method Delete
	r.HandleFunc("/api/v1/fruits/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if id == "" {
			sendResponse(http.StatusBadRequest, "Bad Request, Data Id Params Is Null", nil, w)
			return
		}

		fruit, err := repository.GetFruit(db, id)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error, delete fruit return err", err.Error(), w)
			return
		}

		if fruit.Id == 0 {
			if err != nil {
				sendResponse(http.StatusBadRequest, "bad request", nil, w)
			}
		}

		err = repository.Delete(db, id)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error, delete fruit return err", err.Error(), w)
			return
		}

		sendResponse(http.StatusCreated, "Success Delete", nil, w)
		return
	}).Methods(http.MethodDelete)

	http.HandleFunc("/api/v1/fruits", func(w http.ResponseWriter, r *http.Request) {

		//Get Method
		// if r.Method == http.MethodGet{

		// }

		//Post Method
		// if r.Method == http.MethodPost{

		// }

		//Put Method
		// if r.Method == http.MethodPut{

		// }

		//Delete Method
		// if r.Method == http.MethodDelete{

		// }
	})

	port := "8000"
	fmt.Println("Server Run On Port", port)
	http.ListenAndServe(":"+port, r)
}
