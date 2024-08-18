package entities

type Document struct {
	ID                     string
	Text                   string
	RephrasedText          string
	RephrasedTextEmbedding string
	CreatedBy              string
	TimestampField
}
