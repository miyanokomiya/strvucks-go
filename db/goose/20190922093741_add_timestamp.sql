-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS ' BEGIN NEW.updated_at = now(); RETURN NEW; END; ' language 'plpgsql';

ALTER TABLE users ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE users ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

ALTER TABLE summaries ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE summaries ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON summaries
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

ALTER TABLE permissions ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE permissions ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON permissions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

ALTER TABLE webhook_events ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE webhook_events ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON webhook_events
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users DROP COLUMN created_at;
ALTER TABLE users DROP COLUMN updated_at;
ALTER TABLE summaries DROP COLUMN created_at;
ALTER TABLE summaries DROP COLUMN updated_at;
ALTER TABLE permissions DROP COLUMN created_at;
ALTER TABLE permissions DROP COLUMN updated_at;
ALTER TABLE webhook_events DROP COLUMN created_at;
ALTER TABLE webhook_events DROP COLUMN updated_at;
DROP FUNCTION trigger_set_timestamp() CASCADE;
