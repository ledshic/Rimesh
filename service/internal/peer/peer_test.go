package peer_test

import (
	"testing"
	"time"

	"github.com/ledshic/rimesh/service/internal/peer"
)

func TestNewPeer(t *testing.T) {
	p := peer.NewPeer("id-1", "Alice")
	if p.ID != "id-1" {
		t.Errorf("ID: got %q, want %q", p.ID, "id-1")
	}
	if p.Alias != "Alice" {
		t.Errorf("Alias: got %q, want %q", p.Alias, "Alice")
	}
	if p.GetState() != peer.StateDisconnected {
		t.Errorf("initial state: got %v, want Disconnected", p.GetState())
	}
}

func TestPeerStateTransitions(t *testing.T) {
	p := peer.NewPeer("id-2", "Bob")
	p.SetState(peer.StateConnecting)
	if p.GetState() != peer.StateConnecting {
		t.Fatalf("expected Connecting, got %v", p.GetState())
	}
	p.SetState(peer.StateConnected)
	if p.GetState() != peer.StateConnected {
		t.Fatalf("expected Connected, got %v", p.GetState())
	}
}

func TestPeerUpdateLastSeen(t *testing.T) {
	p := peer.NewPeer("id-3", "Carol")
	before := time.Now()
	p.UpdateLastSeen()
	after := time.Now()
	if p.LastSeen.Before(before) || p.LastSeen.After(after) {
		t.Errorf("LastSeen %v not in [%v, %v]", p.LastSeen, before, after)
	}
}

func TestManagerAddGetRemove(t *testing.T) {
	m := peer.NewManager()

	p1 := peer.NewPeer("id-A", "Alice")
	p2 := peer.NewPeer("id-B", "Bob")
	m.Add(p1)
	m.Add(p2)

	if got := m.Get("id-A"); got != p1 {
		t.Errorf("Get(id-A): got %v, want %v", got, p1)
	}
	if all := m.All(); len(all) != 2 {
		t.Errorf("All(): got %d peers, want 2", len(all))
	}

	m.Remove("id-A")
	if got := m.Get("id-A"); got != nil {
		t.Errorf("after Remove: Get(id-A) should be nil, got %v", got)
	}
	if all := m.All(); len(all) != 1 {
		t.Errorf("after Remove: All() got %d peers, want 1", len(all))
	}
}

func TestManagerGetMissing(t *testing.T) {
	m := peer.NewManager()
	if got := m.Get("nonexistent"); got != nil {
		t.Errorf("Get for missing peer should return nil, got %v", got)
	}
}

func TestStateString(t *testing.T) {
	cases := []struct {
		s    peer.State
		want string
	}{
		{peer.StateDisconnected, "disconnected"},
		{peer.StateConnecting, "connecting"},
		{peer.StateConnected, "connected"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("State(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}
