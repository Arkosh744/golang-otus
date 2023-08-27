-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id         VARCHAR(50) PRIMARY KEY,
    name       TEXT NOT NULL ,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS events
(
    id          VARCHAR(36) PRIMARY KEY,
    user_id     VARCHAR(50) NOT NULL,
    title       TEXT        NOT NULL,
    description TEXT,
    start_at   TIMESTAMP   NOT NULL,
    end_at     TIMESTAMP   NOT NULL,
    notify_at   TIMESTAMP   NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS events;
-- +goose StatementBegin
-- +goose StatementEnd