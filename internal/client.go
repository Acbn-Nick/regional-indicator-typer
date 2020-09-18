package client

import (
	"io/ioutil"
	"runtime"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/getlantern/systray"
	"github.com/lxn/win"
	hook "github.com/robotn/gohook"
	log "github.com/sirupsen/logrus"

	"github.com/Acbn-Nick/regional-indicator-typer/internal/keycode"
)

type Client struct {
	on   bool
	conf *Config
	Done chan interface{}
}

func New() *Client {
	c := NewConfig()

	if err := c.loadConfig(); err != nil {
		log.Infof("failed to load config")
	}

	return &Client{on: false, conf: c}
}

func (c *Client) Start() {
	log.Info("starting client")

	go systray.Run(c.tray, c.onExit)

	go c.hooks()
}

func (c *Client) tray() {
	ico, err := ioutil.ReadFile("../assets/favicon.ico")
	if err != nil {
		log.Fatal("error loading systray icon ", err.Error())
	} else {
		// Add delay to fix issue with systray.AddMenuItem() in goroutines on Windows.
		time.Sleep(time.Second)
		systray.SetIcon(ico)
	}

	systray.SetTitle("Regional Indicator Typer")
	systray.SetTooltip("Regional Indicator Typer")
	quit := systray.AddMenuItem("Quit", "Quit")

	<-quit.ClickedCh

	log.Info("terminating")

	hook.End()
	c.Done <- nil
}

func (c *Client) onExit() {
	c.Done <- nil
}

func (c *Client) hooks() {
	var (
		alphabet = []string{
			"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
			"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v",
			"w", "x", "y", "z",
		}
		emoji = []string{
			"🇦", "🇧", "🇨", "🇩", "🇪", "🇫", "🇬", "🇭", "🇮", "🇯", "🇰",
			"🇱", "🇲", "🇳", "🇴", "🇵", "🇶", "🇷", "🇸", "🇹", "🇺", "🇻",
			"🇼", "🇽", "🇾", "🇿",
		}
		aToE = make(map[string]string)
		wk   = make([]win.KEYBD_INPUT, 4)
	)

	log.Info("starting hooks")

	runtime.LockOSThread()

	sc, err := keycode.Parse(c.conf.Shortcut)
	if err != nil {
		log.Fatalf("failed parsing shortcut: %s", err.Error())
	}

	log.Infof("using shortcut for toggle: %+v", sc)

	hook.Register(hook.KeyDown, sc, func(e hook.Event) {
		log.Infof("emoji mode is now: %+v", !c.on)
		c.on = !c.on
	})

	for i, abc := range alphabet {
		aToE[abc] = emoji[i]
	}

	for i := range alphabet {
		hook.Register(hook.KeyDown, alphabet[i:i+1], func(e hook.Event) {
			if !c.on {
				return
			}

			r, ok := aToE[strings.ToLower(string(e.Keychar))]

			if !ok {
				return
			}

			ru := []rune(r)
			r1, r2 := utf16.EncodeRune(ru[0])

			wk[0] = win.KEYBD_INPUT{
				Type: win.INPUT_KEYBOARD,
				Ki: win.KEYBDINPUT{
					WVk:         0,
					WScan:       uint16(r1),
					DwFlags:     win.KEYEVENTF_UNICODE,
					Time:        0,
					DwExtraInfo: 0,
				},
			}

			wk[1] = win.KEYBD_INPUT{
				Type: win.INPUT_KEYBOARD,
				Ki: win.KEYBDINPUT{
					WVk:         0,
					WScan:       uint16(r2),
					DwFlags:     win.KEYEVENTF_UNICODE,
					Time:        0,
					DwExtraInfo: 0,
				},
			}

			wk[2] = win.KEYBD_INPUT{
				Type: win.INPUT_KEYBOARD,
				Ki: win.KEYBDINPUT{
					WVk:         0,
					WScan:       uint16(r1),
					DwFlags:     win.KEYEVENTF_UNICODE | win.KEYEVENTF_KEYUP,
					Time:        0,
					DwExtraInfo: 0,
				},
			}

			wk[3] = win.KEYBD_INPUT{
				Type: win.INPUT_KEYBOARD,
				Ki: win.KEYBDINPUT{
					WVk:         0,
					WScan:       uint16(r2),
					DwFlags:     win.KEYEVENTF_UNICODE | win.KEYEVENTF_KEYUP,
					Time:        0,
					DwExtraInfo: 0,
				},
			}

			win.SendInput(4, unsafe.Pointer(&wk[0]), int32(unsafe.Sizeof(wk[0])))
			time.Sleep(150 * time.Millisecond)

		})
	}

	s := hook.Start()
	<-hook.Process(s)
}
