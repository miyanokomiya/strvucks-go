-- +goose Up
-- SQL in this section is executed when the migration is applied.
DELETE FROM webhook_events as t1
  WHERE EXISTS(
    SELECT * FROM webhook_events as t2
    WHERE t2.object_id = t1.object_id
    AND t2.ctid > t1.ctid
  );
ALTER TABLE webhook_events ADD constraint webhook_events_unq_object unique(object_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE webhook_events DROP constraint webhook_events_unq_object;
