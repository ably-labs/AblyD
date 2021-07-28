package libablyd

type AblyDInstanceStartMessage struct {
    MessageID string `json:"messageId"`
    ServerID  string `json:"serverID"`
    Args []   string `json:"args"`
}
