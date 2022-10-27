package repository

import (
	entity "HolyCrusade/internal/entity/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *pgx.Conn
}

func (ur UserRepository) Init(db *pgx.Conn) UserRepository {
	ur.db = db

	return ur
}

func (ur *UserRepository) GetByID(ctx context.Context, ID int) (entity.User, error) {
	var u entity.User

	err := ur.db.QueryRow(
		ctx,
		"SELECT id, chat_id FROM users WHERE user_id = $1",
		ID,
	).Scan(&u.ID, &u.ChatID)

	return u, err
}

func (ur *UserRepository) GetByChatId(ctx context.Context, chatID int64) (entity.User, error) {
	var u entity.User

	err := ur.db.QueryRow(
		ctx,
		"SELECT id, chat_id FROM users WHERE chat_id = $1",
		chatID,
	).Scan(&u.ID, &u.ChatID)

	return u, err
}

func (ur *UserRepository) Insert(ctx context.Context, u entity.User) (int, error) {
	var id int

	err := ur.db.QueryRow(
		ctx,
		"INSERT INTO users (chat_id) values ($1) RETURNING id",
		u.ChatID,
	).Scan(&id)

	return id, err
}

func (ur *UserRepository) Delete(ctx context.Context, ID int) error {
	_, err := ur.db.Exec(ctx,
		"DELETE FROM users WHERE id = $1", ID)

	return err
}
