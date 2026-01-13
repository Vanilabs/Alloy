package socket

import (
	"sync"

	"context"

	"time"

	"github.com/gofiber/contrib/websocket"
)

type ConnectionManager struct {
	Connections map[string]*websocket.Conn
	mu          sync.RWMutex
	SocketTracker *SocketTracker
	DMRooms map[string]map[string]bool
}


func NewManager(tracker *SocketTracker) *ConnectionManager {
	return &ConnectionManager{
		Connections: make(map[string]*websocket.Conn),
		SocketTracker: tracker,
		DMRooms:       make(map[string]map[string]bool),
	}
}

func (m *ConnectionManager) Add(socketID string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Connections[socketID] = conn
}

func (m *ConnectionManager) Remove(socketID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Connections, socketID)
}

func (m *ConnectionManager) StartHeartbeat(
	ctx context.Context,
	socketID string,
	conn *websocket.Conn,
) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(
					websocket.PingMessage,
					[]byte("ping"),
					time.Now().Add(5*time.Second),
				); err != nil {
					m.CleanupDeadSocket(socketID)
					return
				}

				m.SocketTracker.RDB.Expire(
					ctx,
					"socket:"+socketID,
					30*time.Second,
				)

			case <-ctx.Done():
				return
			}
		}
	}()
}

func (m *ConnectionManager) CleanupDeadSocket(socketID string) {
	m.mu.Lock()
	conn, exists := m.Connections[socketID]
	if exists {
		_ = conn.Close()
		delete(m.Connections, socketID)
	}
	m.mu.Unlock()

	_ = m.SocketTracker.RemoveSocket(socketID)
}


func (m *ConnectionManager) BroadcastToConnections(message []byte, socketIDs []string) error {

	for _, socketID := range socketIDs {
		m.mu.RLock()
		conn, exists := m.Connections[socketID]
		m.mu.RUnlock()

		if exists {
			// conn.WriteMessage(1, message)
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			// socket is dead → cleanup
			m.CleanupDeadSocket(socketID)
}
		}
	}

	return nil
}


func (m *ConnectionManager) BroadcastToSelf(originSocketID, sender string, message []byte) error {
	receiverSockets, err := m.SocketTracker.GetUserSocketsExcept(sender, originSocketID)
	if err != nil {
		return err
	}

	return m.BroadcastToConnections(message, receiverSockets)
}


func (m *ConnectionManager) BroadcastToAnotherUser(receiver string, message []byte) error {

	receiverSockets, err := m.SocketTracker.GetUserSockets(receiver)
	if err != nil {
		return err
	}

	return m.BroadcastToConnections(message, receiverSockets)
}


