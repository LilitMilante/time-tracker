package tracker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
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
INSERT INTO users (id, passport_series, passport_number, surname, name, patronymic, address, created_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

	_, err := r.db.Exec(ctx, q, u.ID, u.PassportSeries, u.PassportNumber, u.Surname, u.Name, u.Patronymic, u.Address, u.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, updUser UpdateUser) error {
	setSQL, args := setUpdateProductSQL(updUser)

	q := fmt.Sprintf("UPDATE users %s WHERE id = '%s'", setSQL, updUser.ID)

	_, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func setUpdateProductSQL(updUser UpdateUser) (string, []any) {
	var cols []string
	var args []any

	if updUser.PassportSeries != nil {
		args = append(args, *updUser.PassportSeries)
		cols = append(cols, fmt.Sprintf("pasport_series = $%d", len(args)))
	}
	if updUser.PassportNumber != nil {
		args = append(args, *updUser.PassportNumber)
		cols = append(cols, fmt.Sprintf("pasport_number = $%d", len(args)))
	}
	if updUser.Surname != nil {
		args = append(args, *updUser.Surname)
		cols = append(cols, fmt.Sprintf("surname = $%d", len(args)))
	}
	if updUser.Name != nil {
		args = append(args, *updUser.Name)
		cols = append(cols, fmt.Sprintf("name = $%d", len(args)))
	}
	if updUser.Patronymic != nil {
		args = append(args, *updUser.Patronymic)
		cols = append(cols, fmt.Sprintf("patronymic = $%d", len(args)))
	}
	if updUser.Address != nil {
		args = append(args, *updUser.Address)
		cols = append(cols, fmt.Sprintf("address = $%d", len(args)))
	}

	return "SET " + strings.Join(cols, ", "), args
}

func (r *Repository) UserByID(ctx context.Context, id uuid.UUID) (u User, err error) {
	q := `
SELECT id, passport_series, passport_number, surname, name, patronymic, address 
FROM users WHERE id = $1 AND deleted_at ISNULL
`

	err = r.db.QueryRow(ctx, q, id).Scan(
		&u.ID,
		&u.PassportSeries,
		&u.PassportNumber,
		&u.Surname,
		&u.Name,
		&u.Patronymic,
		&u.Address,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return u, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID, deletedAt time.Time) error {
	q := `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at ISNULL`

	res, err := r.db.Exec(ctx, q, deletedAt, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) StartWork(ctx context.Context, wh WorkHours) error {
	q := `
INSERT INTO work_hours (user_id, task_id, started_at)
VALUES ($1, $2, $3)
`

	_, err := r.db.Exec(ctx, q, wh.UserID, wh.TaskID, wh.StartedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FinishWork(ctx context.Context, wh WorkHours) error {
	q := `
UPDATE work_hours
SET finished_at = $1, spend_time_sec = $2
WHERE user_id = $3 AND task_id = $4 AND finished_at ISNULL
`

	_, err := r.db.Exec(ctx, q, wh.FinishedAt, wh.SpendTimeSec, wh.UserID, wh.TaskID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) NotFinishedWorkHours(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (wh WorkHours, err error) {
	q := `SELECT user_id, task_id, started_at, finished_at, spend_time_sec
FROM work_hours
WHERE user_id = $1 AND task_id = $2 AND finished_at ISNULL`

	err = r.db.QueryRow(ctx, q, userID, taskID).Scan(&wh.UserID, &wh.TaskID, &wh.StartedAt, &wh.FinishedAt, &wh.SpendTimeSec)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return WorkHours{}, ErrNotFound
		}
		return WorkHours{}, err
	}

	return wh, nil
}

func (r *Repository) TaskSpendTimesByUser(ctx context.Context, id uuid.UUID) ([]TaskSpendTime, error) {
	q := `
SELECT task_id, SUM(spend_time_sec) FROM work_hours 
WHERE user_id = $1 AND finished_at IS NOT NULL
GROUP BY task_id
`

	rows, err := r.db.Query(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskSpendTimes []TaskSpendTime

	for rows.Next() {
		var taskSpendTime TaskSpendTime

		err = rows.Scan(
			&taskSpendTime.TaskID,
			&taskSpendTime.SpendTimeSec,
		)
		if err != nil {
			return nil, err
		}

		taskSpendTime.UserID = id

		taskSpendTimes = append(taskSpendTimes, taskSpendTime)
	}

	return taskSpendTimes, rows.Err()
}
