# Harmony HID
Share your input devices between multiple machines

This is a work in progress and is in no way functional

## TODO (in no particular order)
- [x] handle active client switching
- [x] websocet server needs a total rewrite
- [x] release all peers on tripple alt
- [x] show/hide cursor as focus moves between peers
- [x] send screen config to server on connect
- [ ] restart peer discovery on connection lost
- [ ] add more logging
- [ ] create ui for arranging screens
- [ ] config file for common settings
- [ ] clean up my shitty code

## Known Bugs
- [ ] input events dont get sent to focused peer for some reason
- [ ] if peer connects too fast after server startup transition zones are assigend incorrectly
- [ ] peer disconnect is not handled properly
- [ ] release events are not always sent to all peers
- [ ] connecting peer sometimes freezes if the server is stopped

## Credits
[github.com/foolin/goview](github.com/foolin/goview)  
[github.com/gin-gonic/gin](github.com/gin-gonic/gin)  
[github.com/gorilla/websocket](github.com/gorilla/websocket)  
[github.com/holoplot/go-evdev](github.com/holoplot/go-evdev)  
[github.com/jezek/xgb](github.com/jezek/xgb)  
[github.com/jezek/xgbutil](github.com/jezek/xgbutil)  
[github.com/vmihailenco/msgpack/v5](github.com/vmihailenco/msgpack/v5)  
