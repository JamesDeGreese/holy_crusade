package repository

import (
	entity "HolyCrusade/internal/entity/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type CityRepository struct {
	db *pgx.Conn
}

func (cr CityRepository) Init(db *pgx.Conn) CityRepository {
	cr.db = db
	return cr
}

func (cr *CityRepository) GetByID(ctx context.Context, ID int) (entity.City, error) {
	var c entity.City

	err := cr.db.QueryRow(
		ctx,
		"SELECT id, user_id, name, rating FROM cities WHERE user_id = $1",
		ID,
	).Scan(&c.ID, &c.UserID, &c.Name, &c.Rating)

	return c, err
}

func (cr *CityRepository) GetByUserID(ctx context.Context, userID int) (entity.City, error) {
	var c entity.City

	err := cr.db.QueryRow(
		ctx,
		"SELECT id, user_id, name, rating FROM cities WHERE user_id = $1",
		userID,
	).Scan(&c.ID, &c.UserID, &c.Name, &c.Rating)

	return c, err
}

func (cr *CityRepository) Insert(ctx context.Context, c entity.City) (int, error) {
	var id int

	err := cr.db.QueryRow(
		ctx,
		"INSERT INTO cities (user_id, name, rating) values ($1, $2, $3) RETURNING id",
		c.UserID,
		c.Name,
		c.Rating,
	).Scan(&id)

	return id, err
}

func (cr *CityRepository) Delete(ctx context.Context, ID int) error {
	_, err := cr.db.Exec(ctx,
		"DELETE FROM cities WHERE id = $1", ID)

	return err
}
