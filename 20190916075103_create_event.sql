-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE webhook_events (
  id                 serial,
  aspect_type        varchar(500),
  event_time         integer,
  object_id          integer,
  object_type        varchar(500),
  owner_id           integer,
  subscription_id    integer,
  PRIMARY KEY (id)
);

CREATE INDEX ON webhook_events (owner_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE webhook_events;
