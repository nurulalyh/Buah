package models

type Fruit struct {
	Id    int     `json: "id"`
	Name  string  `json: "name"`
	Price float64 `json: "price"`
}

func (f Fruit) Validate() error {
	cErr := []string{}
	if f.Name == "" {
		cErr = append(cErr, "name cannot null")
	}
	if f.Price == 0.0 {
		cErr = append(cErr, "price cannot null")
	}

	if len(cErr) > 0 {
		return NewBadRequest(cErr)
	}

	return nil
}

func (f Fruit) Exist() error {
	if f.Id == 0 {
		return NewNotFound("fruit not found")
	}
	return nil
}

type Response struct {
	Code    int         `json: "code"`
	Message string      `json: "message"`
	Data    interface{} `json: "data"`
}
