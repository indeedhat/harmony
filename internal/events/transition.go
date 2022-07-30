package events

import "github.com/indeedhat/harmony/internal/screens"

// TransitionZoneAssigned will be sent to clients on connect and whenever
// the global screen arrangement is updated, it is used to pass the new details
// of their transition zones
type TransitionZoneAssigned []screens.TransitionZone

// Marshal FocusRecieved struct into a byte array for sending via websocket
func (ev TransitionZoneAssigned) Marshal() ([]byte, error) {
	return marshalEvent(ev, MsgTypeTrasitionAssigned)
}

// String gives the string name of the event type
func (tza *TransitionZoneAssigned) String() string {
	return "TransitionZoneAssigned"
}

var _ WsMessage = (*TransitionZoneAssigned)(nil)
