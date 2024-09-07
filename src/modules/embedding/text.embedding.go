package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const OLLAMA_EMBEDDING_URL = "http://ollama.local/api/embeddings"
const NOMICAI_EMBEDDING_URL = "https://api-atlas.nomic.ai/v1/embedding/text"

func floatArrayToStringArray(arr []float64) string {
	// Create a slice to hold the string representations of the integers
	stringSlice := make([]string, len(arr))

	// Convert each float to a string
	for i, num := range arr {
		stringSlice[i] = strconv.FormatFloat(num, 'f', -1, 64)
	}

	// Join the slice into a single string with commas and spaces
	result := "[" + strings.Join(stringSlice, ", ") + "]"

	return result
}

func ollamaGetTextEmbedding(word string) (string, error) {
	// Prepare the request payload
	payload := map[string]string{
		"model":  "nomic-embed-text",
		"prompt": word,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Make the POST request
	resp, err := http.Post(OLLAMA_EMBEDDING_URL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result map[string][]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	embedding, ok := result["embedding"]
	if !ok {
		return "", fmt.Errorf("embedding not found in response")
	}

	return floatArrayToStringArray(embedding), nil
}

func nomicGetTextEmbedding(apiKey string, word string) (string, error) {
	// Create the request payload
	payload := map[string]interface{}{
		"model": "nomic-embed-text-v1.5",
		"texts": []string{word},
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", NOMICAI_EMBEDDING_URL, bytes.NewBuffer(jsonPayload))
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

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// Extract the embedding from the response
	embeddings, ok := result["embeddings"].([]interface{})
	if !ok || len(embeddings) == 0 {
		return "", fmt.Errorf("no embeddings found in response")
	}

	innerEmbeddings, ok := embeddings[0].([]interface{})
	if !ok {
		return "", fmt.Errorf("invalid embedding format")
	}

	// Convert []interface{} to []float64
	floatEmbeddings := make([]float64, len(innerEmbeddings))
	for i, v := range innerEmbeddings {
		num, ok := v.(float64)
		if !ok {
			return "", fmt.Errorf("embedding contains non-float64 value")
		}
		floatEmbeddings[i] = num
	}

	return floatArrayToStringArray(floatEmbeddings), nil
}

func GetTextEmbedding(env string, apiKey string, word string) (string, error) {
	if env == "DEV" {
		return ollamaGetTextEmbedding(word)
	}

	return nomicGetTextEmbedding(apiKey, word)
}
