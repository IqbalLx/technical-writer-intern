package document

import (
	"context"

	"github.com/IqbalLx/technical-writer-intern/src/entities"
	"github.com/IqbalLx/technical-writer-intern/src/modules/embedding"
	"github.com/IqbalLx/technical-writer-intern/src/modules/llm"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DoInsertNewDocument(
	ctx context.Context, db *pgxpool.Pool,
	env string, embedderApiKey string, llmApiKey string,
	text string, createdBy string,
) error {
	paraphrasedWord, err := llm.GetLLMParaphrasedWord(env, llmApiKey, text)
	if err != nil {
		return err
	}
	paraphrasedWordEmbedding, err := embedding.GetTextEmbedding(env, embedderApiKey, paraphrasedWord)
	if err != nil {
		return err
	}

	newDocument := entities.Document{
		Text:                   text,
		RephrasedText:          paraphrasedWord,
		RephrasedTextEmbedding: paraphrasedWordEmbedding,
		CreatedBy:              createdBy,
	}

	return insertNewDocument(ctx, db, newDocument)
}

func DoGetSimilarDocuments(
	ctx context.Context, db *pgxpool.Pool, env string, embedderApiKey string, text string,
) ([]entities.Document, error) {
	textEmbedding, err := embedding.GetTextEmbedding(env, embedderApiKey, text)
	if err != nil {
		return []entities.Document{}, err
	}

	const MAX_COS_DISTANCE = 0.3
	similarDocs, err := querySimilarDocument(ctx, db, textEmbedding, MAX_COS_DISTANCE)
	if err != nil {
		return []entities.Document{}, err
	}

	return similarDocs, nil
}
