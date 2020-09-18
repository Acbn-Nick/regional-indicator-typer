package main

import (
	log "github.com/sirupsen/logrus"

	client "github.com/Acbn-Nick/regional-indicator-typer/internal"
)

func main() {
	log.Info("launching regional-indicator-typer")
	c := client.New()
	c.Start()
	<-c.Done
}
