-- +goose Up
-- +goose StatementBegin
CREATE TABLE events_users
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    event_id   INTEGER                     NOT NULL,
    CONSTRAINT fk_events_users_events
        FOREIGN KEY (event_id)
            REFERENCES events (id),
    user_id    INTEGER                     NOT NULL,
    CONSTRAINT fk_events_users_users
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events_users;
-- +goose StatementEnd
