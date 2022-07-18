package device

import (
	"errors"
	"fmt"

	"github.com/indeedhat/harmony/internal/common"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xinerama"
	"github.com/jezek/xgb/xproto"
)

var _ Vdu = (*X11Vdu)(nil)

// X11Vdu provides common x11 display intergrations
type X11Vdu struct {
	xcon   *xgb.Conn
	window xproto.Window
}

// NewVdu creates a new Vdu isntance,
// in this case for x11
func NewVdu() (Vdu, error) {
	con, err := xgb.NewConn()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to x server: %w", err)
	}

	if err := xinerama.Init(con); err != nil {
		return nil, fmt.Errorf("failed to init xinerama: %w", err)
	}

	setup := xproto.Setup(con)
	if setup == nil {
		return nil, errors.New("failed to setup xproto")
	}

	return X11Vdu{
		xcon:   con,
		window: setup.DefaultScreen(con).Root,
	}, nil
}

// Close the connection to xserver
func (x11 X11Vdu) Close() error {
	x11.xcon.Close()
	return nil
}

// DisplayBounds gets the bounds of the currently connected displays
func (x11 X11Vdu) DisplayBounds() ([]DisplayBounds, error) {
	screens, err := xinerama.QueryScreens(x11.xcon).Reply()
	if err != nil {
		return nil, fmt.Errorf("failed to query screens: %w", err)
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

// CursorPos gets the current coords of the cursor
func (x11 X11Vdu) CursorPos() (*Cursor, error) {
	resp, err := xproto.QueryPointer(x11.xcon, x11.window).Reply()
	if err != nil {
		return nil, fmt.Errorf("failed to query cursor pos: %w", err)
	}

	return &Cursor{
		X: int(resp.RootX),
		Y: int(resp.RootY),
	}, nil
}

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
