package wschat

// IInstance represent the chat application
type IInstance interface {
	BroadcastChannel(request chan []byte)
}
