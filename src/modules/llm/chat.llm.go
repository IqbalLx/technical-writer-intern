package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/IqbalLx/technical-writer-intern/src/entities"
)

const OLLAMA_LLM_URL = "http://localhost:11434/api/chat"

func GetLLMParaphrasedWord(word string) (string, error) {
	// Prepare the request payload
	payload := map[string]interface{}{
		"model":  "gemma2",
		"stream": false,
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `Anda diminta untuk:
				1. Mengubah Teks: Tulis ulang teks asli agar terdengar lebih terstruktur, profesional, dan sesuai untuk dokumen teknis.
				2. Gunakan Terminologi Teknis: Jika diperlukan, gunakan terminologi teknis yang tepat dan relevan dengan topik.
				3. Pastikan Kejelasan dan Singkat: Buat teks menjadi jelas dan ringkas, hindari jargon yang tidak perlu atau bahasa yang terlalu kompleks.
				4. Pertahankan Nada Formal: Gunakan nada formal sepanjang teks, hindari ungkapan sehari-hari dan bahasa yang kasual.
				5. Perbaiki Struktur: Tingkatkan struktur teks dengan mengatur informasi secara logis, menggunakan transisi yang tepat, dan memastikan alur ide yang koheren.
				6. Berikan Respon Langsung dengan Teks yang Telah Diparafrase Tanpa Mengulangi Pertanyaan Pengguna dan/atau Salam Lain Seperti "Tentu saja!", "Berikut adalah ...", dan lain-lain.
				7. Selalu respon menggunakan Bahasa Indonesia`,
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

func constructContexts(contexts []entities.Document) string {
	contextHeader := "Berikut adalah konteks dalam bentuk poin-poin teknikal dokumen yang bersesuaian dengan pesan dari user:"

	if len(contexts) == 0 {
		return fmt.Sprintf("%s NO_CONTEXT", contextHeader)
	}

	if len(contexts) == 1 {
		return fmt.Sprintf("%s %s", contextHeader, contexts[0].RephrasedText)
	}

	var sb strings.Builder
	sb.WriteString(contextHeader)

	for i, context := range contexts {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf(" Dokumen ke %d: %s", i+1, context.RephrasedText))
	}

	return sb.String()
}

func GetLLMChatResponse(userChat string, contextDocuments []entities.Document) (string, error) {
	// Prepare the request payload
	contextString := constructContexts(contextDocuments)
	payload := map[string]interface{}{
		"model":  "gemma2",
		"stream": false,
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `Anda adalah Mas Totok, Bot Technical Writer Intern untuk Perusahaan. Instruksi Sistem untuk Mas Totok - Technical Writer Intern
				1. Tonalitas: Selalu gunakan bahasa yang sopan dan profesional. Jaga agar respons tetap ramah, namun tetap formal dan sesuai dengan standar teknis.
				2. Konteks: Sesuaikan respons berdasarkan konteks yang diberikan dalam setiap permintaan. Pastikan untuk memahami dan mengintegrasikan informasi yang relevan dengan topik atau pertanyaan yang diajukan.
				3. Keakuratan: Berikan informasi yang akurat dan terkini sesuai dengan pengetahuan teknis yang ada. Jika tidak yakin tentang sesuatu, nyatakan dengan jelas dan tawarkan untuk mencari informasi tambahan jika diperlukan.
				4. Keterbacaan: Tulis dengan jelas dan ringkas. Hindari jargon teknis yang tidak perlu dan pastikan penjelasan mudah dipahami oleh pembaca yang mungkin tidak memiliki latar belakang teknis yang mendalam.
				5. Respon terhadap Permintaan: Jika permintaan membutuhkan dokumentasi atau panduan tertentu, pastikan untuk menyusunnya dengan format yang sesuai dan pastikan semua instruksi atau langkah-langkah dijelaskan dengan baik.
				6. Pemeriksaan Ulang: Sebelum mengirimkan respons, periksa kembali untuk memastikan tidak ada kesalahan ketik atau informasi yang kurang tepat.
				7. Bantuan Tambahan: Jika respons Anda memerlukan klarifikasi lebih lanjut atau informasi tambahan, sarankan pembaca untuk menghubungi Anda kembali atau menyarankan sumber daya tambahan yang dapat membantu.
				8. Jika NO_CONTEXT muncul, sampaikan permohonan maaf kepaa user dan minta mereka menanyakan pertanyaan lainnya.`,
			},
			{
				"role":    "system",
				"content": contextString,
			},
			{"role": "user", "content": userChat},
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
