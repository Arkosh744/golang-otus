-- +goose Up
CREATE TABLE IF NOT EXISTS events
(
    id            VARCHAR(36) PRIMARY KEY,
    user_id       VARCHAR(50) NOT NULL,
    title         TEXT        NOT NULL,
    description   TEXT,
    starts_at     TIMESTAMP   NOT NULL,
    ends_at       TIMESTAMP   NOT NULL,
    notify_before BIGINT,
    created_at    TIMESTAMP   NOT NULL DEFAULT NOW()
);
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS events;
-- +goose StatementBegin
-- +goose StatementEnd