package repository

import (
	entity "HolyCrusade/internal/entity/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type CityRepository struct {
	DB *pgx.Conn
}

func (r *CityRepository) GetByID(ID int) (entity.City, error) {
	var c entity.City

	err := r.DB.QueryRow(
		context.Background(),
		"SELECT id, user_id, name, rating FROM cities WHERE user_id = $1",
		ID,
	).Scan(&c.ID, &c.UserID, &c.Name, &c.Rating)

	return c, err
}

func (r *CityRepository) GetByUserID(userID int) (entity.City, error) {
	var c entity.City

	err := r.DB.QueryRow(
		context.Background(),
		"SELECT id, user_id, name, rating FROM cities WHERE user_id = $1",
		userID,
	).Scan(&c.ID, &c.UserID, &c.Name, &c.Rating)

	return c, err
}

func (r *CityRepository) Insert(c entity.City) (int, error) {
	var id int

	err := r.DB.QueryRow(
		context.Background(),
		"INSERT INTO cities (user_id, name, rating) values ($1, $2, $3) RETURNING id",
		c.UserID,
		c.Name,
		c.Rating,
	).Scan(&id)

	return id, err
}

func (r *CityRepository) Delete(ID int) error {
	_, err := r.DB.Exec(context.Background(),
		"DELETE FROM cities WHERE id = $1", ID)

	return err
}
