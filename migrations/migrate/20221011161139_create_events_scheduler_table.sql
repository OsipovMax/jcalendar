-- +goose Up
-- +goose StatementBegin
CREATE TABLE events_scheduler
(
    id               SERIAL PRIMARY KEY,
    created_at       TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_at       TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    begin_occurrence TIMESTAMP                   NOT NULL,
    end_occurrence   TIMESTAMP,
    ending_mode      TEXT                        NOT NULL,
    interval_val     INTEGER                     NOT NULL,
    shift            TEXT                        NOT NULL CHECK (shift IN ('DAILY', 'WEEKLY', 'MONTHLY', 'YEARLY')),
    is_regular       BOOLEAN,
    scheduler_type   TEXT                        NOT NULL,
    event_id         INTEGER                     NOT NULL,
    CONSTRAINT fk_events_scheduler_events
        FOREIGN KEY (event_id)
            REFERENCES events (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events_scheduler;
-- +goose StatementEnd
