-- +goose Up
ALTER TABLE milestones ALTER COLUMN event_date TYPE DATE;

-- +goose Down
ALTER TABLE milestones ALTER COLUMN event_date TYPE TIMESTAMPTZ;
