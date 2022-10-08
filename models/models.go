package models

type Fruit struct {
	Id    int     `json: "id"`
	Name  string  `json: "name"`
	Price float64 `json: "price"`
}

type Response struct {
	Code    int         `json: "code"`
	Message string      `json: "message"`
	Data    interface{} `json: "data"`
}
