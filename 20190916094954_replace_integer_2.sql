-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users ALTER COLUMN athlete_id TYPE bigint;
ALTER TABLE summaries ALTER COLUMN athlete_id TYPE bigint;
ALTER TABLE webhook_events ALTER COLUMN owner_id TYPE bigint;
ALTER TABLE permissions ALTER COLUMN athlete_id TYPE bigint;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users ALTER COLUMN athlete_id TYPE integer;
ALTER TABLE summaries ALTER COLUMN athlete_id TYPE integer;
ALTER TABLE webhook_events ALTER COLUMN owner_id TYPE  integer;
ALTER TABLE permissions ALTER COLUMN athlete_id TYPE integer;
