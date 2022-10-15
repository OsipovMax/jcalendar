-- +goose Up
-- +goose StatementBegin

CREATE TABLE events
(
    id            SERIAL PRIMARY KEY,
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    "from"        TIMESTAMP                   NOT NULL,
    till          TIMESTAMP                   NOT NULL,
    creator_id    INTEGER                     NOT NULL,
    CONSTRAINT fk_events_users
        FOREIGN KEY (creator_id)
            REFERENCES users (id),
    details       TEXT                        NOT NULL,
    schedule_rule TEXT,
    is_private    BOOLEAN                     NOT NULL,
    is_repeat     BOOLEAN                     NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
