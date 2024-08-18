-- +goose Up
-- +goose StatementBegin
CREATE TABLE documents(
    id SERIAL PRIMARY KEY,
    "text" TEXT NOT NULL,
    rephrased_text TEXT NOT NULL,
    rephrased_text_embedding VECTOR(768), -- Nomic AI embedding dims
    created_by TEXT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE INDEX ON documents USING hnsw (rephrased_text_embedding vector_cosine_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE documents;
-- +goose StatementEnd
