package service

import (
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

// Create
func (svc *Service) Create(fruit models.Fruit) error {
	if err := svc.repository.CreateFruits(fruit); err != nil {
		return models.NewInternalServerError(err.Error())
	}
	return nil
}

// Update
func (svc *Service) Update(newFruit models.Fruit) error {
	fruit, err := svc.repository.GetFruit(newFruit.Id)
	if err != nil {
		return models.NewInternalServerError(err.Error())
	}

	err = fruit.Exist()
	if err != nil {
		return err
	}

	fruit.Name = newFruit.Name
	fruit.Price = newFruit.Price

	if err = svc.repository.UpdateFruits(fruit); err != nil {
		return err
	}
	return nil
}

// List
func (svc *Service) List() ([]models.Fruit, error) {
	fruits, err := svc.repository.GetFruits()
	if err != nil {
		return nil, err
	}

	return fruits, nil
}

// Get
func (svc *Service) Get(id int) (models.Fruit, error) {
	fruit, err := svc.repository.GetFruit(id)
	if err != nil {
		return models.Fruit{}, models.NewInternalServerError(err.Error())
	}

	if err = fruit.Exist(); err != nil {
		return models.Fruit{}, err
	}

	return fruit, nil
}

// Delete
func (svc *Service) Delete(id int) error {
	fruit, err := svc.repository.GetFruit(id)
	if err != nil {
		return models.NewInternalServerError(err.Error())
	}

	if err = fruit.Exist(); err != nil {
		return err
	}

	if err = svc.repository.Delete(fruit.Id); err != nil {
		return models.NewInternalServerError(err.Error())
	}

	return nil
}
