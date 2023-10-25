package webservice

import (
	"sync"
)

type PushChannel struct {
	clients      map[*PushClient]bool
	clientsLock  sync.RWMutex
	addClient    chan *PushClient
	removeClient chan *PushClient
	broadcast    chan interface{}
	history      []interface{}
	historyLock  sync.RWMutex
}

func NewPushChannel(broadcastChan chan interface{}) *PushChannel {
	return &PushChannel{
		clients:      make(map[*PushClient]bool),
		addClient:    make(chan *PushClient),
		removeClient: make(chan *PushClient),
		broadcast:    broadcastChan,
		history:      make([]interface{}, 0, 10000),
	}
}

func (c *PushChannel) Run() {
	for {
		select {
		case client := <-c.addClient:
			c.AddClient(client)
		case client := <-c.removeClient:
			c.RemoveClient(client)
		case message := <-c.broadcast:
			c.Broadcast(message)
		}
	}
}

func (c *PushChannel) AddClient(client *PushClient) {
	// Ensure the client is not already registered
	c.clientsLock.RLock()
	if _, ok := c.clients[client]; ok {
		Logger.Warn().Str("client", client.ID).Msg("Websocket client already registered")
		c.clientsLock.RUnlock()
		return
	}
	c.clientsLock.RUnlock()

	// Register the new client
	c.clientsLock.Lock()
	c.clients[client] = true
	c.clientsLock.Unlock()

	Logger.Info().Str("client", client.ID).Msg("Websocket client registered")

	// Send over the channels history to the client, to get the frontend into the correct state
	c.historyLock.RLock()
	for _, msg := range c.history {
		client.tx <- msg
	}
	c.historyLock.RUnlock()

	Logger.Info().Str("client", client.ID).Msg("Websocket client received history")
}

func (c *PushChannel) RemoveClient(client *PushClient) {
	c.clientsLock.Lock()
	if _, ok := c.clients[client]; ok {
		delete(c.clients, client)
		close(client.tx)

		Logger.Info().Str("id", client.ID).Msg("Websocket client unregistered")
	}
	c.clientsLock.Unlock()
}

func (c *PushChannel) Broadcast(message interface{}) {
	c.historyLock.Lock()
	c.history = append(c.history, message)
	c.historyLock.Unlock()

	// Distribute message to all connected clients
	c.clientsLock.RLock()
	clients := c.clients
	for client := range clients {
		select {
		case client.tx <- message:
		default:
			close(client.tx)
		}
	}
	c.clientsLock.RUnlock()
}
