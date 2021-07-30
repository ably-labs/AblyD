package libablyd

type AblyDProcessStartMessage struct {
    MessageID string `json:"messageId"`
    ServerID  string `json:"serverID"`
    Args []   string `json:"args"`
}
