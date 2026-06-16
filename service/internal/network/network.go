// Package network implements the virtual LAN network layer for Rimesh.
// It manages virtual IP allocation, peer routing and signalling.
package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// VirtualSubnet is the subnet used for virtual IP allocation.
	VirtualSubnet = "10.88.0.0/16"

	// SignalPort is the UDP port used for peer signalling.
	SignalPort = 7788

	// HeartbeatInterval is how often keep-alive heartbeats are sent.
	HeartbeatInterval = 10 * time.Second
)

// Message types for the signalling protocol.
const (
	MsgJoin      = "join"
	MsgLeave     = "leave"
	MsgHeartbeat = "heartbeat"
	MsgPeerList  = "peer_list"
)

// Message is the wire format for signalling messages.
type Message struct {
	Type    string          `json:"type"`
	PeerID  string          `json:"peer_id"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// PeerInfo is the payload for join/peer_list messages.
type PeerInfo struct {
	ID        string `json:"id"`
	Alias     string `json:"alias"`
	VirtualIP string `json:"virtual_ip"`
	PublicIP  string `json:"public_ip"`
}

// Network manages the virtual LAN state for the local node.
type Network struct {
	mu sync.RWMutex

	id        string
	alias     string
	virtualIP net.IP
	subnet    *net.IPNet
	conn      *net.UDPConn

	peers    map[string]*PeerInfo
	stopCh   chan struct{}
	doneCh   chan struct{}
	onChange func()
}

// Config holds the configuration for starting a Network.
type Config struct {
	Alias    string
	OnChange func() // called when peer list changes
}

// New creates a new Network instance.
func New(cfg Config) (*Network, error) {
	_, subnet, err := net.ParseCIDR(VirtualSubnet)
	if err != nil {
		return nil, fmt.Errorf("parse virtual subnet: %w", err)
	}

	id := uuid.New().String()
	vip := allocateVirtualIP(subnet, id)

	return &Network{
		id:        id,
		alias:     cfg.Alias,
		virtualIP: vip,
		subnet:    subnet,
		peers:     make(map[string]*PeerInfo),
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
		onChange:  cfg.OnChange,
	}, nil
}

// Start begins listening for UDP signalling messages and announces presence.
func (n *Network) Start() error {
	addr := &net.UDPAddr{Port: SignalPort}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("listen UDP: %w", err)
	}
	n.conn = conn

	go n.receiveLoop()
	go n.heartbeatLoop()

	log.Printf("[network] started: id=%s alias=%s vip=%s", n.id, n.alias, n.virtualIP)
	return nil
}

// Stop gracefully shuts down the network.
func (n *Network) Stop() {
	close(n.stopCh)
	if n.conn != nil {
		_ = n.conn.Close()
	}
	<-n.doneCh
	log.Println("[network] stopped")
}

// ID returns the local node ID.
func (n *Network) ID() string { return n.id }

// VirtualIP returns the assigned virtual IP address.
func (n *Network) VirtualIP() net.IP { return n.virtualIP }

// Peers returns a snapshot of currently known peers.
func (n *Network) Peers() []PeerInfo {
	n.mu.RLock()
	defer n.mu.RUnlock()
	out := make([]PeerInfo, 0, len(n.peers))
	for _, p := range n.peers {
		out = append(out, *p)
	}
	return out
}

// Announce broadcasts a join message to the subnet.
func (n *Network) Announce() error {
	info := PeerInfo{
		ID:        n.id,
		Alias:     n.alias,
		VirtualIP: n.virtualIP.String(),
	}
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}
	msg := Message{Type: MsgJoin, PeerID: n.id, Payload: payload}
	return n.broadcast(msg)
}

func (n *Network) receiveLoop() {
	defer close(n.doneCh)
	buf := make([]byte, 4096)
	for {
		nr, _, err := n.conn.ReadFromUDP(buf)
		if err != nil {
			select {
			case <-n.stopCh:
				return
			default:
				log.Printf("[network] read error: %v", err)
				continue
			}
		}
		n.handleMessage(buf[:nr])
	}
}

func (n *Network) heartbeatLoop() {
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			msg := Message{Type: MsgHeartbeat, PeerID: n.id}
			if err := n.broadcast(msg); err != nil {
				log.Printf("[network] heartbeat error: %v", err)
			}
		case <-n.stopCh:
			return
		}
	}
}

func (n *Network) handleMessage(data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}
	if msg.PeerID == n.id {
		return // ignore own messages
	}

	switch msg.Type {
	case MsgJoin:
		var info PeerInfo
		if err := json.Unmarshal(msg.Payload, &info); err != nil {
			return
		}
		n.mu.Lock()
		n.peers[info.ID] = &info
		n.mu.Unlock()
		log.Printf("[network] peer joined: %s (%s)", info.Alias, info.VirtualIP)
		n.notifyChange()

	case MsgLeave:
		n.mu.Lock()
		delete(n.peers, msg.PeerID)
		n.mu.Unlock()
		log.Printf("[network] peer left: %s", msg.PeerID)
		n.notifyChange()

	case MsgHeartbeat:
		// peer is alive; no state change needed here
	}
}

func (n *Network) broadcast(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	bcastAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: SignalPort,
	}
	_, err = n.conn.WriteToUDP(data, bcastAddr)
	return err
}

func (n *Network) notifyChange() {
	if n.onChange != nil {
		n.onChange()
	}
}

// allocateVirtualIP deterministically assigns a virtual IP within the subnet
// based on the node ID to avoid conflicts on the same local network.
func allocateVirtualIP(subnet *net.IPNet, id string) net.IP {
	// Use a simple hash of the id bytes to pick an offset in the subnet.
	var hash uint32
	for i, b := range []byte(id) {
		hash ^= uint32(b) << (uint32(i%4) * 8)
	}
	// Clamp to a /24 range inside the /16 to reserve the .0 and .255 addresses.
	offset := (hash % 253) + 1
	ip := make(net.IP, 4)
	copy(ip, subnet.IP.To4())
	ip[2] = byte((hash >> 8) % 254)
	ip[3] = byte(offset)
	return ip
}
