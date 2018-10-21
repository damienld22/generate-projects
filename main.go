package main

import (
	"github.com/damienld22/generate-projects/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	cmd.Execute()
}
