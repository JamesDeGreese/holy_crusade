package repository

import (
	entity "HolyCrusade/internal/entity/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type BalanceRepository struct {
	db *pgx.Conn
}

func (br BalanceRepository) Init(db *pgx.Conn) BalanceRepository {
	br.db = db
	return br
}

func (br *BalanceRepository) GetByCityID(ctx context.Context, cityID int) (entity.Balance, error) {
	var b entity.Balance

	err := br.db.QueryRow(
		ctx,
		"SELECT id, city_id, gold, population, workers, solders, heroes FROM balance WHERE city_id = $1",
		cityID,
	).Scan(&b.ID, &b.CityID, &b.Gold, &b.Population, &b.Workers, &b.Solders, &b.Heroes)

	return b, err
}

func (br *BalanceRepository) Insert(ctx context.Context, b entity.Balance) (int, error) {
	var id int

	err := br.db.QueryRow(
		ctx,
		"INSERT INTO balance (city_id, gold, population, workers, solders, heroes) values ($1, $2, $3, $4, $5, $6) RETURNING id",
		b.CityID,
		b.Gold,
		b.Population,
		b.Workers,
		b.Solders,
		b.Heroes,
	).Scan(&id)

	return id, err
}

func (br *BalanceRepository) Delete(ctx context.Context, ID int) error {
	_, err := br.db.Exec(ctx,
		"DELETE FROM balance WHERE id = $1", ID)

	return err
}
