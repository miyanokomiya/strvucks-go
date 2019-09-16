-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE webhook_events ALTER COLUMN event_time TYPE bigint;
ALTER TABLE webhook_events ALTER COLUMN object_id TYPE bigint;
ALTER TABLE webhook_events ALTER COLUMN subscription_id TYPE bigint;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE webhook_events ALTER COLUMN event_time TYPE integer;
ALTER TABLE webhook_events ALTER COLUMN object_id TYPE integer;
ALTER TABLE webhook_events ALTER COLUMN subscription_id TYPE integer;
