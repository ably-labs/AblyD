package libablyd

type NewProcessMessage struct {
    MessageID string
    Pid string
    Namespace string
    ChannelPrefix string
}
