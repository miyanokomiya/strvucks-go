-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE users (
  id            serial,
  athlete_id    integer,
  username      varchar(500),
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ON users (athlete_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;
