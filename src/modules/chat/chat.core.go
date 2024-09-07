package chat

import (
	"fmt"
	"strings"

	"github.com/IqbalLx/technical-writer-intern/src/entities"
)

func constructContexts(contexts []entities.Document) string {
	contextHeader := "Berikut adalah konteks dalam bentuk poin-poin teknikal dokumen yang bersesuaian dengan pesan terbaru dari user:"

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
		sb.WriteString(fmt.Sprintf(" Dokumen ke %d dari %s: %s", i+1, context.CreatedBy, context.RephrasedText))
	}

	return sb.String()
}
