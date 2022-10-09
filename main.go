package main

import (
	"fmt"
	"net/http"

	"github.com/nurulalyh/Buah/config"
	"github.com/nurulalyh/Buah/controllers"
	"github.com/nurulalyh/Buah/repository"
	"github.com/nurulalyh/Buah/router"
	"github.com/nurulalyh/Buah/service"
)

func main() {
	con := config.NewConfig()
	db := config.NewDatabase(&con.DB)

	defer db.Close()

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	controller := controllers.NewController(service)

	r := router.NewRouter(controller)

	fmt.Printf("Server Run On Port %s:%s", con.API.BaseUrl, con.API.Port)
	http.ListenAndServe(con.API.BaseUrl+":"+con.API.Port, r)
}
