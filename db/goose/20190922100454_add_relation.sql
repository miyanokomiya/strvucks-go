-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE permissions ADD FOREIGN KEY (athlete_id) REFERENCES users (athlete_id) ON DELETE CASCADE;
ALTER TABLE webhook_events ADD FOREIGN KEY (owner_id) REFERENCES users (athlete_id) ON DELETE CASCADE;

ALTER TABLE summaries DROP CONSTRAINT summaries_athlete_id_fkey;
ALTER TABLE summaries ADD FOREIGN KEY (athlete_id) REFERENCES users (athlete_id) ON DELETE CASCADE;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE permissions DROP CONSTRAINT permissions_athlete_id_fkey;
ALTER TABLE webhook_events DROP CONSTRAINT webhook_events_owner_id_fkey;
ALTER TABLE summaries DROP CONSTRAINT summaries_athlete_id_fkey;
ALTER TABLE summaries ADD FOREIGN KEY (athlete_id) REFERENCES users (athlete_id);
