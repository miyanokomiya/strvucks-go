-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE permissions DROP COLUMN strava_token;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE permissions ADD COLUMN strava_token varchar(500);
