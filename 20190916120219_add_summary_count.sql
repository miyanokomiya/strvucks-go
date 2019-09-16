-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE summaries ADD COLUMN weekly_count integer;
ALTER TABLE summaries ADD COLUMN monthly_count integer;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE summaries DROP COLUMN weekly_count;
ALTER TABLE summaries DROP COLUMN monthly_count;
