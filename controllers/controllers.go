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

func NewController(service *service.Service) *Controller {
	return &Controller{
		service: *service,
	}
}

func (ctrl *Controller) Create(w http.ResponseWriter, r *http.Request) {
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

	err = ctrl.service.Create(fruit)
	if err != nil {
		sendResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), w)
	}

	sendResponse(http.StatusCreated, "Success", nil, w)
}

func (ctrl *Controller) Update(w http.ResponseWriter, r *http.Request) {
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

	dataByte, err := io.ReadAll(r.Body)

	if err != nil {
		sendResponse(http.StatusBadRequest, "bad request", nil, w)
		return
	}
	defer r.Body.Close()

	var fruit models.Fruit
	err = json.Unmarshal(dataByte, &fruit)
	if err != nil {
		sendResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), w)
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

func (ctrl *Controller) List(w http.ResponseWriter, r *http.Request) {
	fruits, err := ctrl.service.List()
	if err != nil {
		sendResponse(http.StatusInternalServerError, "Internal Server Error, Get Fruits", err.Error(), w)
		return
	}
	sendResponse(http.StatusOK, "Success", fruits, w)
}

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
