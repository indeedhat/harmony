package screens

import (
	"sync"

	"github.com/google/uuid"
	"github.com/indeedhat/harmony/internal/common"
)

// TODO: this is a  temporary solution to get the system working
// it will just put new clients to the far right of existing ones
// once i build the ui it will need to be updated to properly arrange peers

type Peer struct {
	UUID     uuid.UUID
	Displays []DisplayBounds
}

type ScreenManager struct {
	Peers []Peer

	mux *sync.Mutex
}

// NewScreenManager sets up a new manager for screen arrangement and transition
func NewScreenManager() *ScreenManager {
	return &ScreenManager{
		mux: &sync.Mutex{},
	}
}

// AddPeer to the screen manager
// this will regenerate all the transition zones between all peers
func (mgr *ScreenManager) AddPeer(id uuid.UUID, displays []DisplayBounds) map[uuid.UUID][]TransitionZone {
	mgr.mux.Lock()
	defer mgr.mux.Unlock()

	mgr.Peers = append(mgr.Peers, Peer{
		UUID:     id,
		Displays: displays,
	})

	return mgr.buildTransitionZones()
}

// RemovePeer from the screen manager
// this will regenerate all the transition zones between all peers
func (mgr *ScreenManager) RemovePeer(uuid uuid.UUID) map[uuid.UUID][]TransitionZone {
	mgr.mux.Lock()
	defer mgr.mux.Unlock()

	for i, peer := range mgr.Peers {
		if peer.UUID != uuid {
			continue
		}

		mgr.Peers = append(mgr.Peers[:i], mgr.Peers[i:]...)
		break
	}

	return mgr.buildTransitionZones()
}

// PeerExists checks if a peer is already being tracked by the manager
func (mgr *ScreenManager) PeerExists(uuid uuid.UUID) bool {
	for _, peer := range mgr.Peers {
		if peer.UUID == uuid {
			return true
		}
	}

	return false
}

func (mgr *ScreenManager) buildTransitionZones() map[uuid.UUID][]TransitionZone {
	zones := make(map[uuid.UUID][]TransitionZone)

	for _, peer := range mgr.Peers {
		zones[peer.UUID] = []TransitionZone{}
	}

	if len(mgr.Peers) < 2 {
		return zones
	}

	for i := 0; i < len(mgr.Peers)-1; i++ {
		var (
			peerA   = mgr.Peers[i]
			peerB   = mgr.Peers[i+1]
			screenA = findRightMostScreen(peerA.Displays)
			screenB = findLeftMostScreen(peerB.Displays)
			height  = min(screenA.Height, screenB.Height)
		)

		zones[peerA.UUID] = append(zones[peerA.UUID], TransitionZone{
			UUID: peerB.UUID,
			Bounds: [2]common.Vector2{
				{
					X: screenA.Position.X + screenA.Width,
					Y: screenA.Position.Y,
				},
				{
					X: screenA.Position.X + screenA.Width,
					Y: screenA.Position.Y + height,
				},
			},
			Direction: Right,
		})
		zones[peerB.UUID] = append(zones[peerB.UUID], TransitionZone{
			UUID: peerA.UUID,
			Bounds: [2]common.Vector2{
				{
					X: screenB.Position.X,
					Y: screenB.Position.Y,
				},
				{
					X: screenB.Position.X,
					Y: screenB.Position.Y + height,
				},
			},
			Direction: Left,
		})
	}

	return zones
}

func findRightMostScreen(displays []DisplayBounds) DisplayBounds {
	var rightMost *DisplayBounds

	for _, display := range displays {
		if rightMost == nil {
			rightMost = &display
		} else if display.Position.X+display.Width > rightMost.Position.X+rightMost.Width {
			rightMost = &display
		}
	}

	return *rightMost
}

func findLeftMostScreen(displays []DisplayBounds) DisplayBounds {
	var leftMost *DisplayBounds

	for _, display := range displays {
		if leftMost == nil {
			leftMost = &display
		} else if display.Position.X < leftMost.Position.X {
			leftMost = &display
		}
	}

	return *leftMost
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
