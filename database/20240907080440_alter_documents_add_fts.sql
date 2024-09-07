-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents 
ADD COlUMN rephrased_text_tsvector TSVECTOR 
GENERATED ALWAYS AS (to_tsvector('indonesian', rephrased_text)) STORED;

CREATE INDEX documents_search_vector_idx ON documents USING gin(rephrased_text_tsvector);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN rephrased_text_tsvector;
-- +goose StatementEnd
