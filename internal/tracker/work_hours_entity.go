package tracker

import (
	"time"

	"github.com/gofrs/uuid"
)

type WorkHours struct {
	UserID       uuid.UUID  `json:"user_id"`
	TaskID       uuid.UUID  `json:"task_id"`
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	SpendTimeSec int        `json:"spend_time_sec"`
}
