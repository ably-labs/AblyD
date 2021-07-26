package libablyd

type AblyDInstanceStartMessage struct {
    MessageID string `json:"messageId"`
    Args []string   `json:"args"`
}