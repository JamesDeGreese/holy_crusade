package repository

import (
	entity "HolyCrusade/internal/entity/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	DB *pgx.Conn
}

func (r *UserRepository) GetByID(ID int) (entity.User, error) {
	var u entity.User

	err := r.DB.QueryRow(
		context.Background(),
		"SELECT (id, token) FROM users WHERE user_id = $1",
		ID,
	).Scan(&u.ID, &u.Token)

	return u, err
}

func (r *UserRepository) GetByToken(token string) (entity.User, error) {
	var u entity.User

	err := r.DB.QueryRow(
		context.Background(),
		"SELECT (id, token) FROM users WHERE token = $1",
		token,
	).Scan(&u.ID, &u.Token)

	return u, err
}

func (r *UserRepository) Insert(u entity.User) (int, error) {
	var id int

	err := r.DB.QueryRow(
		context.Background(),
		"INSERT INTO users (id, token) values ($1, $2)",
		u.ID,
		u.Token,
	).Scan(&id)

	return id, err
}

func (r *UserRepository) Delete(ID int) error {
	_, err := r.DB.Exec(context.Background(),
		"DELETE FROM users WHERE id = $1", ID)

	return err
}
