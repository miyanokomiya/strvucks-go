-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users ADD COLUMN ifttt_key varchar(500);
ALTER TABLE users ADD COLUMN ifttt_message varchar(500);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users DROP COLUMN ifttt_key;
ALTER TABLE users DROP COLUMN ifttt_message;
