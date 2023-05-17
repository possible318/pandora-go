package reqTypes

type ApiRequest struct {
	Prompt          string `json:"prompt"`
	Model           string `json:"model"`
	MessageId       string `json:"message_id"`
	ParentMessageId string `json:"parent_message_id"`
	ConversationId  string `json:"conversation_id"`
	Stream          bool   `json:"stream"`
}

type GenerateTitleRequest struct {
	MessageID string `json:"message_id"`
	Model     string `json:"model"`
}

type ConversationRequest struct {
	Title     string `json:"title"`
	IsVisible bool   `json:"is_visible"`
}

type GoonRequest struct {
	ConversationId  string `json:"conversation_id"`
	ParentMessageId string `json:"parent_message_id"`
	Model           string `json:"model"`
	Stream          bool   `json:"stream"`
}
