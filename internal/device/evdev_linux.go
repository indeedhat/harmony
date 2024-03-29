package device

import (
	"fmt"
	"io/ioutil"
	"path"
	"syscall"
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/indeedhat/harmony/internal/common"
	"github.com/indeedhat/harmony/internal/events"
)

var _ Device = (*EvdevDevice)(nil)

type EvdevDevice struct {
	dev *evdev.InputDevice
}

// Read an event from the device
// this will block until an event happens or the device is closed
func (evd *EvdevDevice) Read() (*events.InputEvent, error) {
	event, err := evd.dev.ReadOne()
	if err != nil {
		return nil, err
	}

	return &events.InputEvent{
		Time:  event.Time,
		Type:  uint16(event.Type),
		Code:  uint16(event.Code),
		Value: event.Value,
	}, nil
}

// Write an event to the device
func (evd *EvdevDevice) Write(event *events.InputEvent) error {
	return evd.dev.WriteOne(&evdev.InputEvent{
		Time:  event.Time,
		Type:  evdev.EvType(event.Type),
		Code:  evdev.EvCode(event.Code),
		Value: event.Value,
	})
}

// Close the device handle
func (evd *EvdevDevice) Close() error {
	return evd.dev.Close()
}

// String returns a string representation of the device
// it conforms to fmt.Stringer
func (evd *EvdevDevice) String() string {
	var (
		key bool
		rel bool
	)
	for _, typ := range evd.dev.CapableTypes() {
		if typ == evdev.EV_KEY {
			key = true
		}
		if typ == evdev.EV_REL {
			rel = true
		}
	}

	name, _ := evd.dev.Name()
	path := evd.dev.Path()

	return fmt.Sprintf(`%s {
    Name: "%s",
    Path: "%s",
    KeyEvents: "%v",
    RelEvents: "%v",
}`, "EventDevice", name, path, key, rel)
}

// Grab exclusive access to the device
func (evd *EvdevDevice) Grab() error {
	return evd.dev.Grab()
}

// Release the exclusive lock on the device
func (evd *EvdevDevice) Release() error {
	return evd.dev.Ungrab()
}

// ID returns a unique identifer for the underlying device
func (evd *EvdevDevice) ID() string {
	return evd.dev.Path()
}

type EvdevDevicePlus struct {
	EvdevDevice
}

var _ Device = (*EvdevDevicePlus)(nil)
var _ DevicePlus = (*EvdevDevicePlus)(nil)

// MoveCursor relative to its current position
func (edp *EvdevDevicePlus) MoveCursor(move common.Vector2) {
	evTime := syscall.NsecToTimeval(int64(time.Now().Nanosecond()))

	// move x
	edp.dev.WriteOne(&evdev.InputEvent{
		Time:  evTime,
		Type:  evdev.EV_REL,
		Code:  evdev.REL_X,
		Value: int32(move.X),
	})
	// move y
	edp.dev.WriteOne(&evdev.InputEvent{
		Time:  evTime,
		Type:  evdev.EV_REL,
		Code:  evdev.REL_Y,
		Value: int32(move.Y),
	})
	// sync
	edp.dev.WriteOne(&evdev.InputEvent{
		Time:  evTime,
		Type:  evdev.EV_SYN,
		Code:  evdev.SYN_REPORT,
		Value: 0,
	})
}

// FindObservableDevices thot have key or rel events
func FindObservableDevices() (observable []Device) {
	basePath := "/dev/input"

	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := path.Join(basePath, file.Name())

		dev, err := evdev.Open(filePath)
		if err != nil {
			continue
		}

		if !isObservable(dev) {
			dev.Close()
			continue
		}

		observable = append(observable, &EvdevDevice{dev})
	}

	return observable
}

func isObservable(dev *evdev.InputDevice) bool {
	for _, typ := range dev.CapableTypes() {
		if typ != evdev.EV_KEY && typ != evdev.EV_REL {
			continue
		}

		events := dev.CapableEvents(typ)
		if len(events) == 0 {
			continue
		}

		for _, event := range events {
			if event == evdev.REL_X || event == evdev.KEY_SPACE {
				return true
			}
		}
	}

	return false
}

// CreateVirtualDevice that has as much mouse/keyboard capability as possible
// This can be used to pipe all recieved events through on a client machine
func CreateVirtualDevice() (DevicePlus, error) {
	dev, err := evdev.CreateDevice("harmony-virt", evdev.InputID{
		BusType: 0x03,
		Vendor:  0x4712,
		Product: 0x0816,
		Version: 1,
	}, map[evdev.EvType][]evdev.EvCode{
		evdev.EV_REL: {
			evdev.REL_X,
			evdev.REL_Y,
			evdev.REL_Z,
			evdev.REL_RX,
			evdev.REL_RY,
			evdev.REL_RZ,
			evdev.REL_HWHEEL,
			evdev.REL_DIAL,
			evdev.REL_MISC,
			evdev.REL_WHEEL,
		},
		evdev.EV_KEY: {
			evdev.KEY_ESC,
			evdev.KEY_1,
			evdev.KEY_2,
			evdev.KEY_3,
			evdev.KEY_4,
			evdev.KEY_5,
			evdev.KEY_6,
			evdev.KEY_7,
			evdev.KEY_8,
			evdev.KEY_9,
			evdev.KEY_0,
			evdev.KEY_MINUS,
			evdev.KEY_EQUAL,
			evdev.KEY_BACKSPACE,
			evdev.KEY_TAB,
			evdev.KEY_Q,
			evdev.KEY_W,
			evdev.KEY_E,
			evdev.KEY_R,
			evdev.KEY_T,
			evdev.KEY_Y,
			evdev.KEY_U,
			evdev.KEY_I,
			evdev.KEY_O,
			evdev.KEY_P,
			evdev.KEY_LEFTBRACE,
			evdev.KEY_RIGHTBRACE,
			evdev.KEY_ENTER,
			evdev.KEY_LEFTCTRL,
			evdev.KEY_A,
			evdev.KEY_S,
			evdev.KEY_D,
			evdev.KEY_F,
			evdev.KEY_G,
			evdev.KEY_H,
			evdev.KEY_J,
			evdev.KEY_K,
			evdev.KEY_L,
			evdev.KEY_SEMICOLON,
			evdev.KEY_APOSTROPHE,
			evdev.KEY_GRAVE,
			evdev.KEY_LEFTSHIFT,
			evdev.KEY_BACKSLASH,
			evdev.KEY_Z,
			evdev.KEY_X,
			evdev.KEY_C,
			evdev.KEY_V,
			evdev.KEY_B,
			evdev.KEY_N,
			evdev.KEY_M,
			evdev.KEY_COMMA,
			evdev.KEY_DOT,
			evdev.KEY_SLASH,
			evdev.KEY_RIGHTSHIFT,
			evdev.KEY_KPASTERISK,
			evdev.KEY_LEFTALT,
			evdev.KEY_SPACE,
			evdev.KEY_CAPSLOCK,
			evdev.KEY_F1,
			evdev.KEY_F2,
			evdev.KEY_F3,
			evdev.KEY_F4,
			evdev.KEY_F5,
			evdev.KEY_F6,
			evdev.KEY_F7,
			evdev.KEY_F8,
			evdev.KEY_F9,
			evdev.KEY_F10,
			evdev.KEY_NUMLOCK,
			evdev.KEY_SCROLLLOCK,
			evdev.KEY_KP7,
			evdev.KEY_KP8,
			evdev.KEY_KP9,
			evdev.KEY_KPMINUS,
			evdev.KEY_KP4,
			evdev.KEY_KP5,
			evdev.KEY_KP6,
			evdev.KEY_KPPLUS,
			evdev.KEY_KP1,
			evdev.KEY_KP2,
			evdev.KEY_KP3,
			evdev.KEY_KP0,
			evdev.KEY_KPDOT,
			evdev.KEY_ZENKAKUHANKAKU,
			evdev.KEY_102ND,
			evdev.KEY_F11,
			evdev.KEY_F12,
			evdev.KEY_RO,
			evdev.KEY_KATAKANA,
			evdev.KEY_HIRAGANA,
			evdev.KEY_HENKAN,
			evdev.KEY_KATAKANAHIRAGANA,
			evdev.KEY_MUHENKAN,
			evdev.KEY_KPJPCOMMA,
			evdev.KEY_KPENTER,
			evdev.KEY_RIGHTCTRL,
			evdev.KEY_KPSLASH,
			evdev.KEY_SYSRQ,
			evdev.KEY_RIGHTALT,
			evdev.KEY_LINEFEED,
			evdev.KEY_HOME,
			evdev.KEY_UP,
			evdev.KEY_PAGEUP,
			evdev.KEY_LEFT,
			evdev.KEY_RIGHT,
			evdev.KEY_END,
			evdev.KEY_DOWN,
			evdev.KEY_PAGEDOWN,
			evdev.KEY_INSERT,
			evdev.KEY_DELETE,
			evdev.KEY_MACRO,
			evdev.KEY_MUTE,
			evdev.KEY_VOLUMEDOWN,
			evdev.KEY_VOLUMEUP,
			evdev.KEY_POWER,
			evdev.KEY_KPEQUAL,
			evdev.KEY_KPPLUSMINUS,
			evdev.KEY_PAUSE,
			evdev.KEY_SCALE,
			evdev.KEY_KPCOMMA,
			evdev.KEY_HANGEUL,
			evdev.KEY_HANGUEL,
			evdev.KEY_HANJA,
			evdev.KEY_YEN,
			evdev.KEY_LEFTMETA,
			evdev.KEY_RIGHTMETA,
			evdev.KEY_COMPOSE,
			evdev.KEY_STOP,
			evdev.KEY_AGAIN,
			evdev.KEY_PROPS,
			evdev.KEY_UNDO,
			evdev.KEY_FRONT,
			evdev.KEY_COPY,
			evdev.KEY_OPEN,
			evdev.KEY_PASTE,
			evdev.KEY_FIND,
			evdev.KEY_CUT,
			evdev.KEY_HELP,
			evdev.KEY_MENU,
			evdev.KEY_CALC,
			evdev.KEY_SETUP,
			evdev.KEY_SLEEP,
			evdev.KEY_WAKEUP,
			evdev.KEY_FILE,
			evdev.KEY_SENDFILE,
			evdev.KEY_DELETEFILE,
			evdev.KEY_XFER,
			evdev.KEY_PROG1,
			evdev.KEY_PROG2,
			evdev.KEY_WWW,
			evdev.KEY_MSDOS,
			evdev.KEY_COFFEE,
			evdev.KEY_SCREENLOCK,
			evdev.KEY_ROTATE_DISPLAY,
			evdev.KEY_DIRECTION,
			evdev.KEY_CYCLEWINDOWS,
			evdev.KEY_MAIL,
			evdev.KEY_BOOKMARKS,
			evdev.KEY_COMPUTER,
			evdev.KEY_BACK,
			evdev.KEY_FORWARD,
			evdev.KEY_CLOSECD,
			evdev.KEY_EJECTCD,
			evdev.KEY_EJECTCLOSECD,
			evdev.KEY_NEXTSONG,
			evdev.KEY_PLAYPAUSE,
			evdev.KEY_PREVIOUSSONG,
			evdev.KEY_STOPCD,
			evdev.KEY_RECORD,
			evdev.KEY_REWIND,
			evdev.KEY_PHONE,
			evdev.KEY_ISO,
			evdev.KEY_CONFIG,
			evdev.KEY_HOMEPAGE,
			evdev.KEY_REFRESH,
			evdev.KEY_EXIT,
			evdev.KEY_MOVE,
			evdev.KEY_EDIT,
			evdev.KEY_SCROLLUP,
			evdev.KEY_SCROLLDOWN,
			evdev.KEY_KPLEFTPAREN,
			evdev.KEY_KPRIGHTPAREN,
			evdev.KEY_NEW,
			evdev.KEY_REDO,
			evdev.KEY_F13,
			evdev.KEY_F14,
			evdev.KEY_F15,
			evdev.KEY_F16,
			evdev.KEY_F17,
			evdev.KEY_F18,
			evdev.KEY_F19,
			evdev.KEY_F20,
			evdev.KEY_F21,
			evdev.KEY_F22,
			evdev.KEY_F23,
			evdev.KEY_F24,
			evdev.KEY_PLAYCD,
			evdev.KEY_PAUSECD,
			evdev.KEY_PROG3,
			evdev.KEY_PROG4,
			evdev.KEY_ALL_APPLICATIONS,
			evdev.KEY_DASHBOARD,
			evdev.KEY_SUSPEND,
			evdev.KEY_CLOSE,
			evdev.KEY_PLAY,
			evdev.KEY_FASTFORWARD,
			evdev.KEY_BASSBOOST,
			evdev.KEY_PRINT,
			evdev.KEY_HP,
			evdev.KEY_CAMERA,
			evdev.KEY_SOUND,
			evdev.KEY_QUESTION,
			evdev.KEY_EMAIL,
			evdev.KEY_CHAT,
			evdev.KEY_SEARCH,
			evdev.KEY_CONNECT,
			evdev.KEY_FINANCE,
			evdev.KEY_SPORT,
			evdev.KEY_SHOP,
			evdev.KEY_ALTERASE,
			evdev.KEY_CANCEL,
			evdev.KEY_BRIGHTNESSDOWN,
			evdev.KEY_BRIGHTNESSUP,
			evdev.KEY_MEDIA,
			evdev.KEY_SWITCHVIDEOMODE,
			evdev.KEY_KBDILLUMTOGGLE,
			evdev.KEY_KBDILLUMDOWN,
			evdev.KEY_KBDILLUMUP,
			evdev.KEY_SEND,
			evdev.KEY_REPLY,
			evdev.KEY_FORWARDMAIL,
			evdev.KEY_SAVE,
			evdev.KEY_DOCUMENTS,
			evdev.KEY_BATTERY,
			evdev.KEY_BLUETOOTH,
			evdev.KEY_WLAN,
			evdev.KEY_UWB,
			evdev.KEY_UNKNOWN,
			evdev.KEY_VIDEO_NEXT,
			evdev.KEY_VIDEO_PREV,
			evdev.KEY_BRIGHTNESS_CYCLE,
			evdev.KEY_BRIGHTNESS_AUTO,
			evdev.KEY_BRIGHTNESS_ZERO,
			evdev.KEY_DISPLAY_OFF,
			evdev.KEY_WWAN,
			evdev.KEY_WIMAX,
			evdev.KEY_RFKILL,
			evdev.KEY_MICMUTE,
			evdev.BTN_MISC,
			evdev.BTN_0,
			evdev.BTN_1,
			evdev.BTN_2,
			evdev.BTN_3,
			evdev.BTN_4,
			evdev.BTN_5,
			evdev.BTN_6,
			evdev.BTN_7,
			evdev.BTN_8,
			evdev.BTN_9,
			evdev.BTN_MOUSE,
			evdev.BTN_LEFT,
			evdev.BTN_RIGHT,
			evdev.BTN_MIDDLE,
			evdev.BTN_SIDE,
			evdev.BTN_EXTRA,
			evdev.BTN_FORWARD,
			evdev.BTN_BACK,
			evdev.BTN_TASK,
			evdev.BTN_JOYSTICK,
			evdev.BTN_TRIGGER,
			evdev.BTN_THUMB,
			evdev.BTN_THUMB2,
			evdev.BTN_TOP,
			evdev.BTN_TOP2,
			evdev.BTN_PINKIE,
			evdev.BTN_BASE,
			evdev.BTN_BASE2,
			evdev.BTN_BASE3,
			evdev.BTN_BASE4,
			evdev.BTN_BASE5,
			evdev.BTN_BASE6,
			evdev.BTN_DEAD,
			evdev.BTN_GAMEPAD,
			evdev.BTN_SOUTH,
			evdev.BTN_A,
			evdev.BTN_EAST,
			evdev.BTN_B,
			evdev.BTN_C,
			evdev.BTN_NORTH,
			evdev.BTN_X,
			evdev.BTN_WEST,
			evdev.BTN_Y,
			evdev.BTN_Z,
			evdev.BTN_TL,
			evdev.BTN_TR,
			evdev.BTN_TL2,
			evdev.BTN_TR2,
			evdev.BTN_SELECT,
			evdev.BTN_START,
			evdev.BTN_MODE,
			evdev.BTN_THUMBL,
			evdev.BTN_THUMBR,
			evdev.BTN_DIGI,
			evdev.BTN_TOOL_PEN,
			evdev.BTN_TOOL_RUBBER,
			evdev.BTN_TOOL_BRUSH,
			evdev.BTN_TOOL_PENCIL,
			evdev.BTN_TOOL_AIRBRUSH,
			evdev.BTN_TOOL_FINGER,
			evdev.BTN_TOOL_MOUSE,
			evdev.BTN_TOOL_LENS,
			evdev.BTN_TOOL_QUINTTAP,
			evdev.BTN_STYLUS3,
			evdev.BTN_TOUCH,
			evdev.BTN_STYLUS,
			evdev.BTN_STYLUS2,
			evdev.BTN_TOOL_DOUBLETAP,
			evdev.BTN_TOOL_TRIPLETAP,
			evdev.BTN_TOOL_QUADTAP,
			evdev.BTN_WHEEL,
			evdev.BTN_GEAR_DOWN,
			evdev.BTN_GEAR_UP,
			evdev.KEY_OK,
			evdev.KEY_SELECT,
			evdev.KEY_GOTO,
			evdev.KEY_CLEAR,
			evdev.KEY_POWER2,
			evdev.KEY_OPTION,
			evdev.KEY_INFO,
			evdev.KEY_TIME,
			evdev.KEY_VENDOR,
			evdev.KEY_ARCHIVE,
			evdev.KEY_PROGRAM,
			evdev.KEY_CHANNEL,
			evdev.KEY_FAVORITES,
			evdev.KEY_EPG,
			evdev.KEY_PVR,
			evdev.KEY_MHP,
			evdev.KEY_LANGUAGE,
			evdev.KEY_TITLE,
			evdev.KEY_SUBTITLE,
			evdev.KEY_ANGLE,
			evdev.KEY_FULL_SCREEN,
			evdev.KEY_ZOOM,
			evdev.KEY_MODE,
			evdev.KEY_KEYBOARD,
			evdev.KEY_ASPECT_RATIO,
			evdev.KEY_SCREEN,
			evdev.KEY_PC,
			evdev.KEY_TV,
			evdev.KEY_TV2,
			evdev.KEY_VCR,
			evdev.KEY_VCR2,
			evdev.KEY_SAT,
			evdev.KEY_SAT2,
			evdev.KEY_CD,
			evdev.KEY_TAPE,
			evdev.KEY_RADIO,
			evdev.KEY_TUNER,
			evdev.KEY_PLAYER,
			evdev.KEY_TEXT,
			evdev.KEY_DVD,
			evdev.KEY_AUX,
			evdev.KEY_MP3,
			evdev.KEY_AUDIO,
			evdev.KEY_VIDEO,
			evdev.KEY_DIRECTORY,
			evdev.KEY_LIST,
			evdev.KEY_MEMO,
			evdev.KEY_CALENDAR,
			evdev.KEY_RED,
			evdev.KEY_GREEN,
			evdev.KEY_YELLOW,
			evdev.KEY_BLUE,
			evdev.KEY_CHANNELUP,
			evdev.KEY_CHANNELDOWN,
			evdev.KEY_FIRST,
			evdev.KEY_LAST,
			evdev.KEY_AB,
			evdev.KEY_NEXT,
			evdev.KEY_RESTART,
			evdev.KEY_SLOW,
			evdev.KEY_SHUFFLE,
			evdev.KEY_BREAK,
			evdev.KEY_PREVIOUS,
			evdev.KEY_DIGITS,
			evdev.KEY_TEEN,
			evdev.KEY_TWEN,
			evdev.KEY_VIDEOPHONE,
			evdev.KEY_GAMES,
			evdev.KEY_ZOOMIN,
			evdev.KEY_ZOOMOUT,
			evdev.KEY_ZOOMRESET,
			evdev.KEY_WORDPROCESSOR,
			evdev.KEY_EDITOR,
			evdev.KEY_SPREADSHEET,
			evdev.KEY_GRAPHICSEDITOR,
			evdev.KEY_PRESENTATION,
			evdev.KEY_DATABASE,
			evdev.KEY_NEWS,
			evdev.KEY_VOICEMAIL,
			evdev.KEY_ADDRESSBOOK,
			evdev.KEY_MESSENGER,
			evdev.KEY_DISPLAYTOGGLE,
			evdev.KEY_BRIGHTNESS_TOGGLE,
			evdev.KEY_SPELLCHECK,
			evdev.KEY_LOGOFF,
			evdev.KEY_DOLLAR,
			evdev.KEY_EURO,
			evdev.KEY_FRAMEBACK,
			evdev.KEY_FRAMEFORWARD,
			evdev.KEY_CONTEXT_MENU,
			evdev.KEY_MEDIA_REPEAT,
			evdev.KEY_10CHANNELSUP,
			evdev.KEY_10CHANNELSDOWN,
			evdev.KEY_IMAGES,
			evdev.KEY_NOTIFICATION_CENTER,
			evdev.KEY_PICKUP_PHONE,
			evdev.KEY_HANGUP_PHONE,
			evdev.KEY_DEL_EOL,
			evdev.KEY_DEL_EOS,
			evdev.KEY_INS_LINE,
			evdev.KEY_DEL_LINE,
			evdev.KEY_FN,
			evdev.KEY_FN_ESC,
			evdev.KEY_FN_F1,
			evdev.KEY_FN_F2,
			evdev.KEY_FN_F3,
			evdev.KEY_FN_F4,
			evdev.KEY_FN_F5,
			evdev.KEY_FN_F6,
			evdev.KEY_FN_F7,
			evdev.KEY_FN_F8,
			evdev.KEY_FN_F9,
			evdev.KEY_FN_F10,
			evdev.KEY_FN_F11,
			evdev.KEY_FN_F12,
			evdev.KEY_FN_1,
			evdev.KEY_FN_2,
			evdev.KEY_FN_D,
			evdev.KEY_FN_E,
			evdev.KEY_FN_F,
			evdev.KEY_FN_S,
			evdev.KEY_FN_B,
			evdev.KEY_FN_RIGHT_SHIFT,
			evdev.KEY_BRL_DOT1,
			evdev.KEY_BRL_DOT2,
			evdev.KEY_BRL_DOT3,
			evdev.KEY_BRL_DOT4,
			evdev.KEY_BRL_DOT5,
			evdev.KEY_BRL_DOT6,
			evdev.KEY_BRL_DOT7,
			evdev.KEY_BRL_DOT8,
			evdev.KEY_BRL_DOT9,
			evdev.KEY_BRL_DOT10,
			evdev.KEY_NUMERIC_0,
			evdev.KEY_NUMERIC_1,
			evdev.KEY_NUMERIC_2,
			evdev.KEY_NUMERIC_3,
			evdev.KEY_NUMERIC_4,
			evdev.KEY_NUMERIC_5,
			evdev.KEY_NUMERIC_6,
			evdev.KEY_NUMERIC_7,
			evdev.KEY_NUMERIC_8,
			evdev.KEY_NUMERIC_9,
			evdev.KEY_NUMERIC_STAR,
			evdev.KEY_NUMERIC_POUND,
			evdev.KEY_NUMERIC_A,
			evdev.KEY_NUMERIC_B,
			evdev.KEY_NUMERIC_C,
			evdev.KEY_NUMERIC_D,
			evdev.KEY_CAMERA_FOCUS,
			evdev.KEY_WPS_BUTTON,
			evdev.KEY_TOUCHPAD_TOGGLE,
			evdev.KEY_TOUCHPAD_ON,
			evdev.KEY_TOUCHPAD_OFF,
			evdev.KEY_CAMERA_ZOOMIN,
			evdev.KEY_CAMERA_ZOOMOUT,
			evdev.KEY_CAMERA_UP,
			evdev.KEY_CAMERA_DOWN,
			evdev.KEY_CAMERA_LEFT,
			evdev.KEY_CAMERA_RIGHT,
			evdev.KEY_ATTENDANT_ON,
			evdev.KEY_ATTENDANT_OFF,
			evdev.KEY_ATTENDANT_TOGGLE,
			evdev.KEY_LIGHTS_TOGGLE,
			evdev.BTN_DPAD_UP,
			evdev.BTN_DPAD_DOWN,
			evdev.BTN_DPAD_LEFT,
			evdev.BTN_DPAD_RIGHT,
			evdev.KEY_ALS_TOGGLE,
			evdev.KEY_ROTATE_LOCK_TOGGLE,
			evdev.KEY_BUTTONCONFIG,
			evdev.KEY_TASKMANAGER,
			evdev.KEY_JOURNAL,
			evdev.KEY_CONTROLPANEL,
			evdev.KEY_APPSELECT,
			evdev.KEY_SCREENSAVER,
			evdev.KEY_VOICECOMMAND,
			evdev.KEY_ASSISTANT,
			evdev.KEY_KBD_LAYOUT_NEXT,
			evdev.KEY_EMOJI_PICKER,
			evdev.KEY_DICTATE,
			evdev.KEY_BRIGHTNESS_MIN,
			evdev.KEY_BRIGHTNESS_MAX,
			evdev.KEY_KBDINPUTASSIST_PREV,
			evdev.KEY_KBDINPUTASSIST_NEXT,
			evdev.KEY_KBDINPUTASSIST_PREVGROUP,
			evdev.KEY_KBDINPUTASSIST_NEXTGROUP,
			evdev.KEY_KBDINPUTASSIST_ACCEPT,
			evdev.KEY_KBDINPUTASSIST_CANCEL,
			evdev.KEY_RIGHT_UP,
			evdev.KEY_RIGHT_DOWN,
			evdev.KEY_LEFT_UP,
			evdev.KEY_LEFT_DOWN,
			evdev.KEY_ROOT_MENU,
			evdev.KEY_MEDIA_TOP_MENU,
			evdev.KEY_NUMERIC_11,
			evdev.KEY_NUMERIC_12,
			evdev.KEY_AUDIO_DESC,
			evdev.KEY_3D_MODE,
			evdev.KEY_NEXT_FAVORITE,
			evdev.KEY_STOP_RECORD,
			evdev.KEY_PAUSE_RECORD,
			evdev.KEY_VOD,
			evdev.KEY_UNMUTE,
			evdev.KEY_FASTREVERSE,
			evdev.KEY_SLOWREVERSE,
			evdev.KEY_DATA,
			evdev.KEY_ONSCREEN_KEYBOARD,
			evdev.KEY_PRIVACY_SCREEN_TOGGLE,
			evdev.KEY_SELECTIVE_SCREENSHOT,
			evdev.KEY_MACRO1,
			evdev.KEY_MACRO2,
			evdev.KEY_MACRO3,
			evdev.KEY_MACRO4,
			evdev.KEY_MACRO5,
			evdev.KEY_MACRO6,
			evdev.KEY_MACRO7,
			evdev.KEY_MACRO8,
			evdev.KEY_MACRO9,
			evdev.KEY_MACRO10,
			evdev.KEY_MACRO11,
			evdev.KEY_MACRO12,
			evdev.KEY_MACRO13,
			evdev.KEY_MACRO14,
			evdev.KEY_MACRO15,
			evdev.KEY_MACRO16,
			evdev.KEY_MACRO17,
			evdev.KEY_MACRO18,
			evdev.KEY_MACRO19,
			evdev.KEY_MACRO20,
			evdev.KEY_MACRO21,
			evdev.KEY_MACRO22,
			evdev.KEY_MACRO23,
			evdev.KEY_MACRO24,
			evdev.KEY_MACRO25,
			evdev.KEY_MACRO26,
			evdev.KEY_MACRO27,
			evdev.KEY_MACRO28,
			evdev.KEY_MACRO29,
			evdev.KEY_MACRO30,
			evdev.KEY_MACRO_RECORD_START,
			evdev.KEY_MACRO_RECORD_STOP,
			evdev.KEY_MACRO_PRESET_CYCLE,
			evdev.KEY_MACRO_PRESET1,
			evdev.KEY_MACRO_PRESET2,
			evdev.KEY_MACRO_PRESET3,
		},
	})

	if err != nil {
		return nil, err
	}

	return &EvdevDevicePlus{EvdevDevice{dev}}, nil
}
