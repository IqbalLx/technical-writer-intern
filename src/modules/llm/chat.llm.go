package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/IqbalLx/technical-writer-intern/src/entities"
)

const OLLAMA_LLM_URL = "http://ollama.local/api/chat"

func GetLLMParaphrasedWord(word string) (string, error) {
	// Prepare the request payload
	payload := map[string]interface{}{
		"model":  "gemma2",
		"stream": false,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": WordParaphraserPrompt,
			},
			{"role": "user", "content": word},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Make the POST request
	resp, err := http.Post(OLLAMA_LLM_URL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	message, ok := result["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("message not found in response or not a string")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content not found in response or not a string")
	}

	return content, nil
}

func GetLLMChatResponse(chatHistories []entities.ChatHistory) (string, error) {
	// prepare message
	messages := []map[string]string{
		{
			"role":    "system",
			"content": MasTotokPrompt,
		},
	}

	for _, chatHistory := range chatHistories {
		messages = append(
			messages,
			map[string]string{
				"role":    chatHistory.Role,
				"content": chatHistory.Content,
			})
	}

	payload := map[string]interface{}{
		"model":    "llama3.1",
		"stream":   false,
		"messages": messages,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Make the POST request
	resp, err := http.Post(OLLAMA_LLM_URL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	message, ok := result["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("message not found in response or not a string")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content not found in response or not a string")
	}

	return content, nil
}
