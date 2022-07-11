package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/holoplot/go-evdev"
	"github.com/indeedhat/harmony/internal/device"
	"github.com/indeedhat/harmony/internal/net"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	serverAddress = "192.168.0.10:8080"
)

func main() {
	newDev, err := device.CreateVirtualDevice()
	if err != nil {
		log.Fatal("failed to clone device: ", err.Error())
	}

	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/ws"}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	msg := net.ClientConnect{
		Hostname: "test-client",
	}
	data, err := msg.Marshal()
	if err != nil {
		log.Fatal("failed to marshal connect message")
	}

	ws.WriteMessage(websocket.BinaryMessage, data)

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			return
		}

		switch data[0] {
		case byte(net.MsgTypeSHidEvent):
			var msg net.ServerHidEvent
			if err := msgpack.Unmarshal(data[2:], &msg); err != nil {
				continue
			}

			newDev.WriteOne(&evdev.InputEvent{
				Time:  msg.Time,
				Type:  evdev.EvType(msg.Type),
				Code:  evdev.EvCode(msg.Code),
				Value: msg.Value,
			})

		default:
			log.Print("invalid message type: ", string(data))
		}
	}
}
