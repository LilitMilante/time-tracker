package tracker

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	PassportSeries int       `json:"passport_series"`
	PassportNumber int       `json:"passport_number"`
	Surname        string    `json:"surname"`
	Name           string    `json:"name"`
	Patronymic     string    `json:"patronymic"`
	Address        string    `json:"address"`
	CreatedAt      time.Time `json:"created_at"`
}

type UpdateUser struct {
	ID             uuid.UUID `json:"id"`
	PassportSeries *int      `json:"passport_series"`
	PassportNumber *int      `json:"passport_number"`
	Surname        *string   `json:"surname"`
	Name           *string   `json:"name"`
	Patronymic     *string   `json:"patronymic"`
	Address        *string   `json:"address"`
}

type UserFilter struct {
	ID             *uuid.UUID
	PassportSeries *int
	PassportNumber *int
	Surname        *string
	Name           *string
	Patronymic     *string
	Address        *string
}
