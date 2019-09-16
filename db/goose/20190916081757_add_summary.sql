-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE summaries (
  id                            serial,
  athlete_id                    integer REFERENCES users(athlete_id),
  month_base_date               date,
  monthly_distance              numeric,
  monthly_moving_time           numeric,
  monthly_total_elevation_gain  numeric,
  monthly_calories              numeric,
  week_base_date                date,
  weekly_distance               numeric,
  weekly_moving_time            numeric,
  weekly_total_elevation_gain   numeric,
  weekly_calories               numeric,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ON summaries (athlete_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE summaries;
