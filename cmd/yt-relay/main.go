package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/mirror-media/yt-relay/cli"
	"github.com/mirror-media/yt-relay/cli/serve"
)

func main() {

	cmds := map[string]*cli.Command{
		"serve": serve.Command,
	}

	err := cli.Start(cmds)
	if err != nil {
		log.Fatal(err)
	}

}
