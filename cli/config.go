package cli

import (
	"flag"

	"github.com/mirror-media/yt-relay/config"
)

type Conf struct {
	Address    string
	ConfigFile string
	Port       int
	CFG        *config.Conf
}

func registerFlags(c *Conf, f *flag.FlagSet) {
	f.StringVar(&c.Address, "address", "127.0.0.1", "Address to bind")
	f.StringVar(&c.ConfigFile, "config", "", "path to the configuration file")
	f.IntVar(&c.Port, "port", 80, "Port to bind")
}
