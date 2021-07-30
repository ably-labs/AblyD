package libablyd

type AblyDProcessState struct {
    ServerID string
    Namespace string

    MaxProcesses int
    Processes map[int]string
}