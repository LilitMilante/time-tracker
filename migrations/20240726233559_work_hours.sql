-- +goose Up
-- +goose StatementBegin
CREATE TABLE work_hours (
    user_id UUID NOT NULL REFERENCES users (id),
    task_id UUID NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ,
    spend_time_sec INTEGER NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE work_hours;
-- +goose StatementEnd
