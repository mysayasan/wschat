package wschat

// instance struct
type instance struct {
	hub *Hub
}

// NewInstance New Modbus Client instance
func NewInstance(hub *Hub) IInstance {
	return &instance{
		hub: hub,
	}
}

// RunBroadcaster to run request by listening to channel
func (a *instance) BroadcastChannel(bcastChan chan []byte) {
	go a.broadcastChannel(bcastChan)
}

func (a *instance) broadcastChannel(bcastChan chan []byte) {
	for {
		if message, ok := <-bcastChan; ok {
			a.Broadcast(message)
		}
	}
}

// Broadcast to clients
func (a *instance) Broadcast(message []byte) {
	for client := range a.hub.Clients {
		client.Send <- message
	}
}
