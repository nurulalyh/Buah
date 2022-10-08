package service

import (
	"errors"

	"github.com/nurulalyh/Buah/models"
	"github.com/nurulalyh/Buah/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

// create
func (svc *Service) Create(fruit models.Fruit) error {
	err := svc.repository.CreateFruits(fruit)
	if err != nil {
		return err
	}
	return nil
}

// update
func (svc *Service) Update(newFruit models.Fruit) error {
	fruit, err := svc.repository.GetFruit(newFruit.Id)
	if err != nil {
		return err
	}

	if fruit.Id == 0 {
		return errors.New("data Fruits Not Found")
	}

	fruit.Name = newFruit.Name
	fruit.Price = newFruit.Price

	err = svc.repository.UpdateFruits(fruit)
	if err != nil {
		return err
	}
	return nil
}

// get all
func (svc *Service) List() ([]models.Fruit, error) {
	fruits, err := svc.repository.GetFruits()
	if err != nil {
		return nil, err
	}
	return fruits, nil
}

// get satu
func (svc *Service) Get(id int) (models.Fruit, error) {
	fruit, err := svc.repository.GetFruit(id)
	if err != nil {
		return models.Fruit{}, err
	}
	return fruit, nil
}

// delete
func (svc *Service) Delete(id int) error {
	fruit, err := svc.repository.GetFruit(id)
	if err != nil {
		return err
	}

	if fruit.Id == 0 {
		return errors.New("data Fruits Not Found")
	}

	err = svc.repository.Delete(fruit.Id)
	if err != nil {
		return err
	}
	return nil
}
