package manager

import (
	"sync"

	"github.com/sakokazuki/simplegrpc/event"
)

type Client struct {
	events      []string
	payloadChan chan event.Payload
}

func (c *Client) ReceivePayload() <-chan event.Payload {
	return c.payloadChan
}

func (c *Client) SetEvents(events []string) {
	c.events = events
}

func NewClient(events []string) Client {
	return Client{
		events:      events,
		payloadChan: make(chan event.Payload, 20),
	}
}

type Clients struct {
	clients map[chan event.Payload]struct{}
	mu      *sync.RWMutex
}

type ClientManager struct {
	clientsTable map[string]Clients
}

func (cm *ClientManager) AddClient(client Client) {
	for _, e := range client.events {
		if cm.clientsTable[e].clients == nil {
			cm.clientsTable[e] = Clients{
				clients: make(map[chan event.Payload]struct{}),
				mu:      &sync.RWMutex{},
			}
		}
		cm.clientsTable[e].mu.Lock()
		cm.clientsTable[e].clients[client.payloadChan] = struct{}{}
		cm.clientsTable[e].mu.Unlock()
	}
}

func (cm *ClientManager) RemoveClient(client Client) {
	for _, e := range client.events {
		clients, ok := cm.clientsTable[e]
		if !ok {
			continue
		}
		clients.mu.Lock()
		delete(clients.clients, client.payloadChan)
		clients.mu.Unlock()
	}
	close(client.payloadChan)
}

func (cm *ClientManager) DeleteEvents(client *Client) {
	for _, e := range client.events {
		clients, ok := cm.clientsTable[e]
		if !ok {
			continue
		}
		clients.mu.Lock()
		delete(clients.clients, client.payloadChan)
		clients.mu.Unlock()
	}
}

func (cm *ClientManager) SendPayload(payload event.Payload) {
	eventname := payload.Meta.Type

	clientsTable, ok := cm.clientsTable[eventname]
	if !ok {
		return
	}
	clientsTable.mu.RLock()
	clients := clientsTable.clients
	clientsTable.mu.RUnlock()

	wg2 := sync.WaitGroup{}
	for client := range clients {
		wg2.Add(1)
		go func(client chan event.Payload) {
			defer wg2.Done()
			sendPayloadSafety(client, payload)

		}(client)
	}
	wg2.Wait()
}

func sendPayloadSafety(client chan event.Payload, payload event.Payload) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	client <- payload
	return
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clientsTable: make(map[string]Clients),
	}
}
