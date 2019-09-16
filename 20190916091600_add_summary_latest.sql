-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE summaries ADD COLUMN latest_distance numeric;
ALTER TABLE summaries ADD COLUMN latest_moving_time numeric;
ALTER TABLE summaries ADD COLUMN latest_total_elevation_gain numeric;
ALTER TABLE summaries ADD COLUMN latest_calories numeric;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE summaries DROP COLUMN latest_distance;
ALTER TABLE summaries DROP COLUMN latest_moving_time;
ALTER TABLE summaries DROP COLUMN latest_total_elevation_gain;
ALTER TABLE summaries DROP COLUMN latest_calories;
