-- +goose Up
-- +goose StatementBegin
CREATE UNLOGGED TABLE chat_histories(
    user_name TEXT,
    "role" TEXT CHECK ("role" in ('system', 'user', 'assistant')),
    content TEXT,
    "timestamp" TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE INDEX chat_histories_idx ON chat_histories (user_name, "timestamp");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chat_histories;
-- +goose StatementEnd
