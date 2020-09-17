package client

import (
	"io/ioutil"
	"time"

	"github.com/getlantern/systray"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	conf *Config
}

func New() *Client {
	c := NewConfig()

	if err := c.loadConfig(); err != nil {
		log.Infof("failed to load config")
	}

	return &Client{conf: c}
}

func (c *Client) Start() {
	log.Info("starting client")

	go c.tray()

}

func (c *Client) tray() {
	ico, err := ioutil.ReadFile("../assets/favicon.ico")
	if err != nil {
		log.Fatal("error loading systray icon ", err.Error())
	}
	time.Sleep(500 * time.Millisecond) // Add 500ms delay to fix issue with systray.AddMenuItem() in goroutines on Windows.
	systray.SetIcon(ico)
	systray.SetTitle("Regional Indicator Typer")
	systray.SetTooltip("Regional Indicator Typer")
	quit := systray.AddMenuItem("Quit", "Quit")
}
