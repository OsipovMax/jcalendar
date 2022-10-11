-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id               SERIAL PRIMARY KEY,
    created_at       TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_at       TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    first_name       TEXT    NOT NULL,
    last_name        TEXT    NOT NULL,
    email            TEXT    NOT NULL,
    time_zone_offset INTEGER NOT NULL,
    hashed_password  TEXT    NOT NULL
);

CREATE UNIQUE INDEX unique_index_users_on_email ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP INDEX unique_index_users_on_email;
-- +goose StatementEnd
