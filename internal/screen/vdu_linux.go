package screen

import (
	"github.com/indeedhat/harmony/internal/common"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xinerama"
)

// GetDisplayBounds for all connected monitors
func GetDisplayBounds(cons ...*xgb.Conn) ([]DisplayBounds, error) {
	var (
		con *xgb.Conn
		err error
	)

	if len(cons) == 0 {
		con, err = common.InitXCon()
		if err != nil {
			return nil, err
		}
		defer con.Close()
	} else {
		con = cons[0]
	}

	xinerama.Init(con)
	screens, err := xinerama.QueryScreens(con).Reply()
	if err != nil {
		return nil, err
	}

	count := int(screens.Number)
	displays := make([]DisplayBounds, 0, count)

	for i := 0; i < count; i++ {
		screen := screens.ScreenInfo[i]
		displays = append(displays, DisplayBounds{
			X:      int(screen.XOrg),
			Y:      int(screen.YOrg),
			Width:  int(screen.Width),
			Height: int(screen.Height),
		})
	}

	return displays, nil
}
