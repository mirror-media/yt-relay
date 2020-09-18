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
	f.StringVar(&c.Address, "address", "0.0.0.0", "Address to bind")
	f.StringVar(&c.ConfigFile, "config", "", "path to the configuration file")
	f.IntVar(&c.Port, "port", 8080, "Port to bind")
}
