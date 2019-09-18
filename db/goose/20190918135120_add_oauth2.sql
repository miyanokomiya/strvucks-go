-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE permissions ADD COLUMN access_token varchar(500);
ALTER TABLE permissions ADD COLUMN token_type varchar(500);
ALTER TABLE permissions ADD COLUMN refresh_token varchar(500);
ALTER TABLE permissions ADD COLUMN expiry bigint;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE permissions DROP COLUMN access_token;
ALTER TABLE permissions DROP COLUMN token_type;
ALTER TABLE permissions DROP COLUMN refresh_token;
ALTER TABLE permissions DROP COLUMN expiry;
