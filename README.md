# Harmony HID
Share your input devices between multiple machines

## Current State
This is very much a work in progress.  

Harmony is currently in a usable if not user friendly state: see TODO list for details

## Quick start 
- clone the repository
```sh
git clone git@github.com:indeedhat/harmony
# alternatively
git clone github.com/indeedhat/harmony
```

- copy the configs/config.toml file to the root of the repo
```sh
cd harmony
cp configs/config.toml .
```
> I recommend that you change the cluster_id from the default but everything else should be fine unless a port conflicts

- build the project
```sh
make build
# alternatively if you dont have make
CGO_ENABLED=0 go build -o . ./...
```

- run harmony
```sh
./harmony-hid
```
> Peer discovery is done automatically so starting up the harmony binary on another machine (so long as it has the 
> same cluster_id) will automatically connect to the cluster
> new clients peers will always be set positioned at the right of the previous peers screen

- move your mouse to the far right of your monitor/multi monitor setup to take control of the next peer

Copy the [config file](https://github.com/indeedhat/harmony/blob/master/configs/config.toml) to the directory harmony 
is installed to

## TODO (in no particular order)
- [x] handle active client switching
- [x] websocet server needs a total rewrite
- [x] release all peers on tripple alt
- [x] show/hide cursor as focus moves between peers
- [x] send screen config to server on connect
- [x] config file for common settings
- [x] write logs to file by default with option to print
- [x] allow multiple clusters to runn independently on the same network (currently all instances will join the same cluster)
- [x] restart peer discovery on connection lost
- [x] handle case where multiple servers are spun up on the same network
- [ ] create ui for arranging screens
- [ ] clean up my shitty code
- [ ] clipboard support
- [ ] drag and drop files?
- [ ] windows support
- [ ] macos support (not sure how im gonna do this as i don't have one)

## Known bugs
currently none known

## Credits
[github.com/foolin/goview](github.com/foolin/goview)  
[github.com/gin-gonic/gin](github.com/gin-gonic/gin)  
[github.com/gorilla/websocket](github.com/gorilla/websocket)  
[github.com/holoplot/go-evdev](github.com/holoplot/go-evdev)  
[github.com/jezek/xgb](github.com/jezek/xgb)  
[github.com/jezek/xgbutil](github.com/jezek/xgbutil)  
[github.com/vmihailenco/msgpack/v5](github.com/vmihailenco/msgpack/v5)  
