-- +goose Up
-- +goose StatementBegin
CREATE TABLE invites
(
    id          SERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    user_id     INTEGER                     NOT NULL,
    CONSTRAINT fk_invites_users
        FOREIGN KEY (user_id)
            REFERENCES users (id),
    event_id    INTEGER                     NOT NULL,
    CONSTRAINT fk_invites_events
        FOREIGN KEY (event_id)
            REFERENCES events (id),
    is_accepted BOOLEAN                     NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE invites;
-- +goose StatementEnd
