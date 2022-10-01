package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"strconv"
	"database/sql"
	_ "github.com/lib/pq"
)

type Fruit struct {
	// Id int `json: "id"`
	Name string `json: "name"`
	Price float64 `json: "price"`
}

type Response struct {
	Code int 		 `json: "code"`
	Message string 	 `json: "message"`
	Data interface{} `json: "data"`
}

// var fruits []Fruit

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

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:rootroot@localhost/fruits?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}

	if err = db.Ping(); err != nil {
		panic(err.Error())
	}

	defer db.Close()

	http.HandleFunc("/api/v1/fruits", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet{
			// fmt.
			rows, err := db.Query("select Name, Price from fruits")
			if err != nil{
				sendResponse(http.StatusBadRequest, "Internal Server Error", nil, w)
			}

			var fruits []Fruit
fmt.Println(rows == nil)
			for rows.Next(){
				var fruit Fruit

				err=rows.Scan(
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

			// dataByte, err := json.Marshal(fruits)
			// if err != nil {
			// 	w.Header().Set("content-type", "application/json")
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	w.Write([]byte("Ada Error" + err.Error()))
			// 	return	
			// }
			// w.Header().Set("content-type", "application/json")
			// w.WriteHeader(http.StatusOK)
			// w.Write(dataByte)
			// return
		}

		if r.Method == http.MethodPost{
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

			// w.Write([]byte("Ini Post atau Create"))
			sendResponse(http.StatusCreated, "Success", nil, w)
			return
		}

		if r.Method == http.MethodPut{
			id := r.URL.Query().Get("id")

			if id == "" {
				sendResponse(http.StatusBadRequest, "Bad Request, Data Id Params Is Null", nil, w)
				return
			}

			//cek id ada atau tidak
			idInt, err := strconv.Atoi(id)

			if err != nil {
				sendResponse(http.StatusInternalServerError,"Internal Server Error, Fail Convert String To Int", nil, w)
				return
			}

			// idInt -= 1
			dataByte, err := io.ReadAll(r.Body)

			if err != nil {
				sendResponse(http.StatusBadRequest, "bad request", nil, w)
			}
			
			defer r.Body.Close()
			var fruit Fruit
			err = json.Unmarshal(dataByte, &fruit)

			if err != nil{
				sendResponse(http.StatusInternalServerError, "Internal Server Error", nil, w)
			}

			_, err = db.Exec("UPDATE fruits SET name=$1, price=$2 WHERE id=$3", fruit.Name, fruit.Price, idInt)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, get fruits", nil, w)
			}

			sendResponse(http.StatusCreated, "Success Update", nil, w)
			return

			// w.Write([]byte("Ini Put"))
		}

		if r.Method == http.MethodDelete{
			id := r.URL.Query().Get("id")

			if id == "" {
				sendResponse(http.StatusBadRequest, "Bad Request, Data Id Params Is Null", nil, w)
				return
			}

			//cek id ada atau tidak
			idInt, err := strconv.Atoi(id)

			if err != nil {
				sendResponse(http.StatusInternalServerError,"Internal Server Error, Fail Convert String To Int", nil, w)
				return
			}

			_, err = db.Exec("DELETE from fruits WHERE id=$1", idInt)
			if err != nil {
				sendResponse(http.StatusInternalServerError, "internal server error, get fruits", nil, w)
			}

			sendResponse(http.StatusCreated, "Success Delete", nil, w)
			return
			// w.Write([]byte("Ini Delete"))
		}

		// w.Write([]byte("Wrong Method"))
	})

	port := "8000"
	fmt.Println("Server Run On Port", port)
	http.ListenAndServe(":"+port, nil)
}