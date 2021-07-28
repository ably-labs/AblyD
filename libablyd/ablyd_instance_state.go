package libablyd

type AblyDInstanceState struct {
    ServerID string
    Namespace string

    MaxInstances int
    Instances map[int]string
}