package repository

import (
	"database/sql"
	"errors"

	"github.com/nurulalyh/Buah/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create
func (r *Repository) CreateFruits(fruit models.Fruit) (err error) {
	_, err = r.db.Exec("insert into fruits(Name,Price) values($1,$2)", fruit.Name, fruit.Price)
	if err != nil {
		return errors.New("create fruit return error : " + err.Error())
	}

	return nil
}

// List
func (r *Repository) GetFruits() (fruits []models.Fruit, err error) {
	rows, err := r.db.Query("select Id, Name, Price from fruits")
	if err != nil {
		return nil, errors.New("get fruit return error : " + err.Error())
	}

	for rows.Next() {
		var fruit models.Fruit

		err = rows.Scan(
			&fruit.Id,
			&fruit.Name,
			&fruit.Price,
		)

		if err != nil {
			return nil, errors.New("get fruit return error, scan data return error  : " + err.Error())
		}
		fruits = append(fruits, fruit)
	}
	return fruits, nil
}

// Get
func (r *Repository) GetFruit(id int) (fruit models.Fruit, err error) {
	rows, err := r.db.Query("select id, Name, Price from fruits where id = $1", id)
	if err != nil {
		return models.Fruit{}, errors.New("get fruit return error : " + err.Error())
	}

	if rows.Next() {
		err = rows.Scan(
			&fruit.Id,
			&fruit.Name,
			&fruit.Price,
		)

		if err != nil {
			return models.Fruit{}, errors.New("get fruit return error, scan data return error  : " + err.Error())
		}

	}
	return fruit, nil
}

// Delete
func (r *Repository) Delete(id int) (err error) {
	_, err = r.db.Exec("delete from fruits WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

// Update
func (r *Repository) UpdateFruits(fruit models.Fruit) (err error) {
	_, err = r.db.Exec("update fruits set name = $2, price = $3 where id = $1", fruit.Id, fruit.Name, fruit.Price)
	if err != nil {
		return err
	}

	return nil
}
