package chat

import (
	"context"
	"math"

	"github.com/IqbalLx/technical-writer-intern/src/entities"
	"github.com/IqbalLx/technical-writer-intern/src/modules/document"
	"github.com/IqbalLx/technical-writer-intern/src/modules/llm"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DoAnswerUserChat(ctx context.Context, db *pgxpool.Pool, userName string, userChat string) (string, error) {
	LIMIT_HISTORIES := 40

	currentHistoriesCount, err := countUserChatHistories(ctx, db, userName)
	if err != nil {
		return "", err
	}
	historiesOffset := int(math.Max(float64(currentHistoriesCount-LIMIT_HISTORIES), 0))

	chatHistories, err := getUserChatHistories(ctx, db, userName, historiesOffset, LIMIT_HISTORIES)
	if err != nil {
		return "", err
	}

	contexts, err := document.DoGetSimilarDocuments(ctx, db, userChat)
	if err != nil {
		return "", err
	}

	contextStr := constructContexts(contexts)
	newChats := []entities.ChatHistory{{Role: "system", Content: contextStr}, {Role: "user", Content: userChat}}
	chatHistories = append(chatHistories, newChats...)

	botResponse, err := llm.GetLLMChatResponse(chatHistories)
	if err != nil {
		return "", err
	}

	newChats = append(newChats, entities.ChatHistory{Role: "assistant", Content: botResponse})

	err = insertNewChatHistories(ctx, db, userName, newChats)
	if err != nil {
		return "", err
	}

	return botResponse, nil
}
