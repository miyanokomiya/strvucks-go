-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE permissions (
  id            serial,
  athlete_id    integer,
  strava_token    varchar(500)
);

CREATE UNIQUE INDEX ON permissions (athlete_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE permissions;
