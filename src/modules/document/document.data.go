package document

import (
	"context"

	"github.com/IqbalLx/technical-writer-intern/src/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leporo/sqlf"
)

func insertNewDocument(ctx context.Context, db *pgxpool.Pool, document entities.Document) error {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		InsertInto("documents").
		NewRow().
		Set("text", document.Text).
		Set("rephrased_text", document.RephrasedText).
		Set("rephrased_text_embedding", document.RephrasedTextEmbedding).
		Set("created_by", document.CreatedBy)

	sql, args := query.String(), query.Args()
	if _, err := db.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func querySimilarDocument(ctx context.Context, db *pgxpool.Pool, embedding string,
	maxDistance float64) ([]entities.Document, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		Select("rephrased_text, created_by").
		From("documents").
		Clause("WHERE").Expr("(rephrased_text_embedding <=> ?) <= ?", embedding, maxDistance).
		Clause("ORDER BY").Expr("rephrased_text_embedding <=> ?", embedding).
		Limit(5)

	sql, args := query.String(), query.Args()
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return []entities.Document{}, err
	}
	docs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.Document, error) {
		var doc entities.Document
		err = row.Scan(
			&doc.RephrasedText,
			&doc.CreatedBy,
		)
		if err != nil {
			return doc, err
		}

		return doc, err
	})

	if err != nil {
		return []entities.Document{}, err
	}

	return docs, nil
}
