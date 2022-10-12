package repository

import (
	entity "HolyCrusade/internal/entity/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type BalanceRepository struct {
	DB *pgx.Conn
}

func (r *BalanceRepository) GetByCityID(cityID int) (entity.Balance, error) {
	var b entity.Balance

	err := r.DB.QueryRow(
		context.Background(),
		"SELECT (id, city_id, gold, population, workers, solders, heroes) FROM balance WHERE city_id = $1",
		cityID,
	).Scan(&b.ID, &b.CityID, &b.Gold, &b.Population, &b.Workers, &b.Solders, &b.Heroes)

	return b, err
}

func (r *BalanceRepository) Insert(b entity.Balance) (int, error) {
	var id int

	err := r.DB.QueryRow(
		context.Background(),
		"INSERT INTO balance (city_id, gold, population, workers, solders, heroes) values ($1, $2, $3, $4, $5, $6)",
		b.CityID,
		b.Gold,
		b.Population,
		b.Workers,
		b.Solders,
		b.Heroes,
	).Scan(&id)

	return id, err
}

func (r *BalanceRepository) Delete(ID int) error {
	_, err := r.DB.Exec(context.Background(),
		"DELETE FROM balance WHERE id = $1", ID)

	return err
}
