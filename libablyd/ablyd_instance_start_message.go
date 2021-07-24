package libablyd

type AblyDInstanceStartMessage struct {
    MessageID string `json:"messageId"`
    Action string   `json:"action"`
    Args []string   `json:"args"`
}