package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nurulalyh/Buah/models"
	"github.com/nurulalyh/Buah/service"
)

type Controller struct {
	service service.Service
}

func sendResponse(code int, message string, data interface{}, w http.ResponseWriter) {
	resp := models.Response{
		Code:    code,
		Data:    data,
		Message: message,
	}
	write(resp, code, w)
}

func sendResponseError(err error, w http.ResponseWriter) {
	if v, ok := err.(*models.Errors); ok {
		write(v, v.Code, w)
		return
	}

	data := models.NewErrors(http.StatusInternalServerError, "internal server error", err.Error())
	write(data, http.StatusInternalServerError, w)
}

func write(data interface{}, code int, w http.ResponseWriter) {
	dataByte, _ := json.Marshal(data)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dataByte)
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		service: *service,
	}
}

// Create
func (ctrl *Controller) Create(w http.ResponseWriter, r *http.Request) {
	dataByte, err := io.ReadAll(r.Body)
	if err != nil {
		sendResponseError(models.NewInternalServerError(err.Error()), w)
		return
	}
	defer r.Body.Close()

	var fruit models.Fruit
	err = json.Unmarshal(dataByte, &fruit)

	if err != nil {
		sendResponseError(models.NewInternalServerError(err.Error()), w)
		return
	}

	if err = fruit.Validate(); err != nil {
		sendResponseError(err, w)
		return
	}

	err = ctrl.service.Create(fruit)
	if err != nil {
		sendResponseError(err, w)
		return
	}

	sendResponse(http.StatusCreated, "Success created data", nil, w)
}

// Update
func (ctrl *Controller) Update(w http.ResponseWriter, r *http.Request) {
	//get query param
	id := mux.Vars(r)["id"]

	if id == "" {
		sendResponseError(models.NewBadRequest("error parameter id"), w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendResponseError(models.NewBadRequest("cannot convert id param from string to int"), w)
		return
	}

	dataByte, err := io.ReadAll(r.Body)

	if err != nil {
		sendResponseError(err, w)
		return
	}
	defer r.Body.Close()

	var fruit models.Fruit
	err = json.Unmarshal(dataByte, &fruit)
	if err != nil {
		sendResponseError(err, w)
		return
	}

	fruit.Id = idInt

	err = ctrl.service.Update(fruit)
	if err != nil {
		sendResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), w)
		return
	}

	sendResponse(http.StatusOK, "Success Update", nil, w)
}

// Delete
func (ctrl *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	//get query param
	id := mux.Vars(r)["id"]

	if id == "" {
		sendResponse(http.StatusBadRequest, "Bad Request, Data Id Params Is Null", nil, w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendResponse(http.StatusBadRequest, "cannot convert id param from string to int", nil, w)
		return
	}

	err = ctrl.service.Delete(idInt)
	if err != nil {
		sendResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), w)
		return
	}

	sendResponse(http.StatusOK, "Success Update", nil, w)
}

// List
func (ctrl *Controller) List(w http.ResponseWriter, r *http.Request) {
	fruits, err := ctrl.service.List()
	if err != nil {
		sendResponse(http.StatusInternalServerError, "Internal Server Error, Get Fruits", err.Error(), w)
		return
	}
	sendResponse(http.StatusOK, "Success", fruits, w)
}

// Get
func (ctrl *Controller) Get(w http.ResponseWriter, r *http.Request) {
	//get query param
	id := mux.Vars(r)["id"]

	if id == "" {
		sendResponse(http.StatusBadRequest, "Bad Request, Data Id Params Is Null", nil, w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendResponse(http.StatusBadRequest, "cannot convert id param from string to int", nil, w)
		return
	}

	fruit, err := ctrl.service.Get(idInt)
	if err != nil {
		sendResponse(http.StatusBadRequest, "internal server error, get fruit", err.Error(), w)
		return
	}

	sendResponse(http.StatusOK, "Success", fruit, w)
}
