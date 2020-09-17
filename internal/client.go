package client

import (
	log "github.com/sirupsen/logrus"
)

type client struct {
}

func New() *client {
	return &client{}
}

func (c *client) Start() {
	log.Info("starting client")

}

func (c *client) tray() {

}
