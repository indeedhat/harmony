package mouse

import (
	"log"
	"sync"

	"github.com/indeedhat/harmony/internal/common"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/mousebind"
	"github.com/jezek/xgbutil/xevent"
)

func GetCursorPos(cons ...*xgb.Conn) (x int, y int, err error) {
	var con *xgb.Conn
	if len(cons) == 0 {
		con, err := common.InitXCon()
		if err != nil {
			return x, y, err
		}
		defer con.Close()
	} else {
		con = cons[0]
	}

	// screens, err := vdu.GetDisplayBounds(con)
	// if err != nil {
	// 	return x, y, err
	// }

	setup := xproto.Setup(con)
	window := setup.DefaultScreen(con).Root
	repl, err := xproto.QueryPointer(con, window).Reply()
	if err != nil {
		return x, y, err
	}

	x = int(repl.RootX)
	y = int(repl.RootY)
	return
}

func DebugMouseEvents() error {
	con, err := xgbutil.NewConn()
	if err != nil {
		return err
	}

	buttonPress := mousebind.ButtonPressFun(func(_ *xgbutil.XUtil, event xevent.ButtonPressEvent) {
		log.Printf("click %#v", event)
	})
	buttonPress.Connect(con, con.RootWin(), "", false, true)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

	return nil
}
