package libablyd

type AblyDInstanceState struct {
    MaxInstances int
    Instances map[int]string
}