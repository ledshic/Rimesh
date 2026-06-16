// Package peer manages P2P peer connections for the Rimesh virtual LAN.
package peer

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// State represents the connection state of a peer.
type State int

const (
	StateDisconnected State = iota
	StateConnecting
	StateConnected
)

func (s State) String() string {
	switch s {
	case StateConnecting:
		return "connecting"
	case StateConnected:
		return "connected"
	default:
		return "disconnected"
	}
}

// Peer represents a remote participant in the virtual LAN.
type Peer struct {
	mu sync.RWMutex

	ID        string
	Alias     string
	VirtualIP net.IP
	PublicIP  net.UDPAddr
	State     State
	LastSeen  time.Time
}

// NewPeer creates a new Peer with the given ID and alias.
func NewPeer(id, alias string) *Peer {
	return &Peer{
		ID:    id,
		Alias: alias,
		State: StateDisconnected,
	}
}

// UpdateLastSeen records the current time as the last seen timestamp.
func (p *Peer) UpdateLastSeen() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LastSeen = time.Now()
}

// SetState updates the peer connection state.
func (p *Peer) SetState(s State) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.State = s
}

// GetState returns the current connection state.
func (p *Peer) GetState() State {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.State
}

// String returns a human-readable representation of the peer.
func (p *Peer) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return fmt.Sprintf("Peer{ID:%s Alias:%s VIP:%s State:%s}", p.ID, p.Alias, p.VirtualIP, p.State)
}

// Manager manages the set of known peers.
type Manager struct {
	mu    sync.RWMutex
	peers map[string]*Peer
}

// NewManager creates a new peer Manager.
func NewManager() *Manager {
	return &Manager{
		peers: make(map[string]*Peer),
	}
}

// Add registers a peer with the manager.
func (m *Manager) Add(p *Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.peers[p.ID] = p
}

// Remove unregisters a peer by ID.
func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.peers, id)
}

// Get retrieves a peer by ID. Returns nil if not found.
func (m *Manager) Get(id string) *Peer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.peers[id]
}

// All returns a snapshot of all peers.
func (m *Manager) All() []*Peer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Peer, 0, len(m.peers))
	for _, p := range m.peers {
		out = append(out, p)
	}
	return out
}
