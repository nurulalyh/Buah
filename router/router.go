package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nurulalyh/Buah/controllers"
)

func NewRouter(ctrl *controllers.Controller) *mux.Router {
	r := mux.NewRouter()

	//List
	r.HandleFunc("/api/v1/fruits", ctrl.List).Methods(http.MethodGet)

	//Method Get
	r.HandleFunc("/api/v1/fruits/{id}", ctrl.Get).Methods(http.MethodGet)

	//Method Put
	r.HandleFunc("/api/v1/fruits/{id}", ctrl.Update).Methods(http.MethodPut)

	//Method Post
	r.HandleFunc("/api/v1/fruits", ctrl.Create).Methods(http.MethodPost)

	//Method Delete
	r.HandleFunc("/api/v1/fruits/{id}", ctrl.Delete).Methods(http.MethodDelete)

	return r
}
