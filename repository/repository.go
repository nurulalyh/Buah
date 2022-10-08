package repository

import (
	"database/sql"

	"github.com/nurulalyh/Buah/models"
)

//Create

func CreateFruits(db *sql.DB, fruit models.Fruit) (err error) {
	_, err = db.Exec("insert into fruits(Name,Price) values($1,$2)", fruit.Name, fruit.Price)
	if err != nil {
		return err
	}

	return nil
}

//Get

func GetFruits(db *sql.DB) (fruits []models.Fruit, err error) {
	rows, err := db.Query("select id, Name, Price from fruits")
	if err != nil {
		return nil, err
	}

	// var fruits []models.Fruit
	// fmt.Println(rows == nil)
	for rows.Next() {
		var fruit models.Fruit

		err = rows.Scan(
			&fruit.Id,
			&fruit.Name,
			&fruit.Price,
		)

		if err != nil {
			return nil, err
		}
		fruits = append(fruits, fruit)
	}
	return fruits, nil
}

// Delete
func GetFruit(db *sql.DB, id string) (fruit models.Fruit, err error) {
	rows, err := db.Query("select id, Name, Price from fruits where id = $1", id)
	if err != nil {
		return models.Fruit{}, err
	}

	// var fruits []models.Fruit
	// fmt.Println(rows == nil)
	if rows.Next() {
		// var fruit models.Fruit

		err = rows.Scan(
			&fruit.Id,
			&fruit.Name,
			&fruit.Price,
		)

		if err != nil {
			return models.Fruit{}, err
		}
		// fruits = append(fruits, fruit)
	}
	return fruit, nil
}

func Delete(db *sql.DB, id string) (err error) {
	_, err = db.Exec("DELETE from fruits WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

// update
func UpdateFruits(db *sql.DB, fruit models.Fruit) (err error) {
	_, err = db.Exec("update fruits set name = $2, price = $3 where id = $1", fruit.Id, fruit.Name, fruit.Price)
	if err != nil {
		return err
	}

	return nil
}
