package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Fruit struct {
	Id int `json: "id"`
	Name string `json: "name"`
	Price float64 `json: "price"`
}

type Response struct {
	Code int 		 `json: "code"`
	Message string 	 `json: "message"`
	Data interface{} `json: "data"`
}

func sendResponse(code int, message string, data interface{}, w http.ResponseWriter) {
	resp := Response {
		Code: 	 code,
		Data: 	 data,
		Message: message,
	}
	dataByte, err := json.Marshal(resp)

	if err != nil {
		resp := Response{
			Code: 	 http.StatusInternalServerError,
			Data: 	 nil,
			Message: "Internal Server Error",
		}
		dataByte, _ = json.Marshal(resp)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dataByte)
}

func Remove(slice []Fruit, s int) []Fruit {
	return append(slice[:1], slice[s+1:]...)
}

var db *sql.DB

func GetFruit (w http.ResponseWriter, r *http.Request) {
	
}

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

	r.HandleFunc("/api/v1/fruits", func (w http.ResponseWriter, r *http.Request)  {
		fmt.Println(1111,db == nil)
		rows, err := db.Query("select id, Name, Price from fruits")
			if err != nil{
				sendResponse(http.StatusBadRequest, "Internal Server Error", nil, w)
				return
			}

			var fruits []Fruit
			fmt.Println(rows == nil)
			for rows.Next(){
				var fruit Fruit

				err=rows.Scan(
					&fruit.Id,
					&fruit.Name,
					&fruit.Price,
				)

				if err != nil {
					sendResponse(http.StatusInternalServerError, "Internal Server", nil, w)
				}
				fruits = append(fruits, fruit)
			}
			sendResponse(http.StatusOK, "Success", fruits, w)
			return
	}).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/fruits/{id}", func (w http.ResponseWriter, r *http.Request) {
		//get query param
		id := mux.Vars(r)["id"]

		if id == "" {
			sendResponse(http.StatusBadRequest, "Bad Request, Data Id Params Is Null", nil, w)
			return
		}

		// idInt, err := strconv.Atoi(id)
		rows, err := db.Query("select id, name, price from fruits where id = $1", id)
		if err != nil {
			sendResponse(http.StatusInternalServerError,"Internal Server Error, Get Fruits", nil, w)
			return
		}

		var fruit Fruit
		if rows.Next() {
			err = rows.Scan(
				&fruit.Id,
				&fruit.Name,
				&fruit.Price,
			)

			if err != nil {
				sendResponse(http.StatusInternalServerError,"Internal Server Error, Scan data return err", nil, w)
				return
			}
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

		var newFruit Fruit
		err = json.Unmarshal(dataByte, &newFruit)
		if err != nil{
			sendResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), w)
			return
		}

		fruit.Name = newFruit.Name
		fruit.Price = newFruit.Price

		_, err = db.Exec("UPDATE fruits SET name=$2, price=$3 WHERE id=$1", fruit.Id, fruit.Name, fruit.Price)
		if err != nil {
			sendResponse(http.StatusInternalServerError, "internal server error, get fruits", err.Error(), w)
			return
		}
		sendResponse(http.StatusOK, "Success Update", nil, w)
		return
	}).Methods(http.MethodPut)

	r.HandleFunc("/api/v1/fruits", func (w http.ResponseWriter, r *http.Request) {
		dataByte, err := io.ReadAll(r.Body)
			if err != nil {
				sendResponse(http.StatusBadRequest, "Bad Request", nil, w)
			}
			defer r.Body.Close()

			var fruit Fruit
			err = json.Unmarshal(dataByte, &fruit)

			if err != nil {
				sendResponse(http.StatusInternalServerError, "Internal Server Error", nil, w)
			}

			_, err = db.Exec("insert into fruits(Name,Price) values($1,$2)", fruit.Name, fruit.Price)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, get fruits", nil, w)
			}
			sendResponse(http.StatusCreated, "Success", nil, w)
			return
	}).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/fruits/{id}", func (w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

			if id == "" {
				sendResponse(http.StatusBadRequest, "Bad Request, Data Id Params Is Null", nil, w)
				return
			}

			rows, err := db.Query("select id, name, price from fruits where id = $1", id)
			if err != nil {
				sendResponse(http.StatusInternalServerError,"Internal Server Error, Get Fruits", nil, w)
				return
			}

			var fruit Fruit
			if rows.Next() {
				err = rows.Scan(
					&fruit.Id,
					&fruit.Name,
					&fruit.Price,
				)

				if err != nil {
					sendResponse(http.StatusInternalServerError,"Internal Server Error, Scan data return err", nil, w)
					return
				}
			}

			if fruit.Id == 0 {
				if err != nil {
					sendResponse(http.StatusBadRequest, "bad request", nil, w)
				}
			}

			_, err = db.Exec("DELETE from fruits WHERE id=$1", fruit.Id)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, delete fruit return err", nil, w)
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