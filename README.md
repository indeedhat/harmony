# Harmony HID
Share your input devices between multiple machines

## Current State
This is very much a work in progress.  
Although it is technically in a working state and does allow sharing HID devices between peers (linux x11 only)
it is lacking a lot of functionality and is still a buggy mess

## TODO (in no particular order)
- [x] handle active client switching
- [x] websocet server needs a total rewrite
- [x] release all peers on tripple alt
- [x] show/hide cursor as focus moves between peers
- [x] send screen config to server on connect
- [ ] restart peer discovery on connection lost
- [ ] handle case where multiple servers are spun up on the same network
- [ ] allow multiple setups to runn independently on the same network (currently all instances will join the same cluster)
- [ ] create ui for arranging screens
- [ ] config file for common settings
- [ ] clean up my shitty code
- [ ] write logs to file by default with option to print

## Known Bugs
- [ ] if peer connects too fast after server startup transition zones are assigend incorrectly
- [ ] peer disconnect is not always handled properly
- [ ] release events are not always sent to all peers

## Credits
[github.com/foolin/goview](github.com/foolin/goview)  
[github.com/gin-gonic/gin](github.com/gin-gonic/gin)  
[github.com/gorilla/websocket](github.com/gorilla/websocket)  
[github.com/holoplot/go-evdev](github.com/holoplot/go-evdev)  
[github.com/jezek/xgb](github.com/jezek/xgb)  
[github.com/jezek/xgbutil](github.com/jezek/xgbutil)  
[github.com/vmihailenco/msgpack/v5](github.com/vmihailenco/msgpack/v5)  
