package chat

import (
	"context"

	"github.com/IqbalLx/technical-writer-intern/src/modules/document"
	"github.com/IqbalLx/technical-writer-intern/src/modules/llm"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DoAnswerUserChat(ctx context.Context, db *pgxpool.Pool, userChat string) (string, error) {
	contexts, err := document.DoGetSimilarDocuments(ctx, db, userChat)
	if err != nil {
		return "", err
	}

	botResponse, err := llm.GetLLMChatResponse(userChat, contexts)
	if err != nil {
		return "", err
	}

	return botResponse, nil
}
