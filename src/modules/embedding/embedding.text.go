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

func GetTextEmbedding(word string) (string, error) {
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
