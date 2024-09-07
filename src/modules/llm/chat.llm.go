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
const GROQ_LLM_URL = "https://api.groq.com/openai/v1/chat/completions"

const OLLAMA_LLAMA_MODEL = "llama3.1"
const OLLAMA_GEMMA_MODEL = "gemma2"
const GROQ_LLAMA_MODEL = "llama-3.1-8b-instant"
const GROQ_GEMMA_MODEL = "gemma2-9b-it"

// common
func prepareLLMChatPayload(provider string, chatHistories []entities.ChatHistory) ([]byte, error) {
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

	modelName := OLLAMA_LLAMA_MODEL
	if provider == "groq" {
		modelName = GROQ_LLAMA_MODEL
	}

	payload := map[string]interface{}{
		"model":    modelName,
		"stream":   false,
		"messages": messages,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	return payloadBytes, nil
}

func prepareLLMParaphrasedWordPayload(provider string, word string) ([]byte, error) {
	modelName := OLLAMA_GEMMA_MODEL
	if provider == "groq" {
		modelName = GROQ_GEMMA_MODEL
	}

	// Prepare the request payload
	payload := map[string]interface{}{
		"model":  modelName,
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
		return []byte{}, fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	return payloadBytes, nil
}

// provider
func useOllama(payloadBytes []byte) (string, error) {
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

func useGroq(apiKey string, payloadBytes []byte) (string, error) {
	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", GROQ_LLM_URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set the required headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse JSON response into a generic map
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	// Extract content from the response
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices found in the response")
	}

	// Access the first choice and its message
	choice := choices[0].(map[string]interface{})
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message format")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content is not a string")
	}

	return content, nil
}

// main
func GetLLMParaphrasedWord(env string, apiKey string, word string) (string, error) {
	provider := "ollama"
	if env == "PROD" {
		provider = "groq"
	}

	payloadBytes, err := prepareLLMParaphrasedWordPayload(provider, word)
	if err != nil {
		return "", err
	}

	if provider == "ollama" {
		return useOllama(payloadBytes)
	}

	return useGroq(apiKey, payloadBytes)
}

func GetLLMChatResponse(env string, apiKey string, chatHistories []entities.ChatHistory) (string, error) {
	provider := "ollama"
	if env == "PROD" {
		provider = "groq"
	}

	payloadBytes, err := prepareLLMChatPayload(provider, chatHistories)
	if err != nil {
		return "", err
	}

	if provider == "ollama" {
		return useOllama(payloadBytes)
	}

	return useGroq(apiKey, payloadBytes)
}
