package tracker

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, u User) error {
	q := `
INSERT INTO users (id, passport_series, passport_number, surname, name, patronymic, address) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

	_, err := r.db.Exec(ctx, q, u.ID, u.PassportSeries, u.PassportNumber, u.Surname, u.Name, u.Patronymic, u.Address)
	if err != nil {
		return err
	}

	return nil
}
