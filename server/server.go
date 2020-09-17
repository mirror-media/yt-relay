package server

import (
	"fmt"

	ytrelay "github.com/mirror-media/yt-relay"
	"github.com/mirror-media/yt-relay/config"
	"github.com/mirror-media/yt-relay/whitelist"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Server struct {
	APIWhitelist ytrelay.APIWhitelist
	conf         *config.Conf
	Engine       *gin.Engine
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}

func (s *Server) Run() error {
	return s.Engine.Run(fmt.Sprintf("%s:%d", s.conf.Address, s.conf.Port))
}

func New(c config.Conf) (*Server, error) {

	engine := gin.Default()

	s := &Server{
		APIWhitelist: &whitelist.API{
			Whitelist: c.Whitelists,
		},
		conf:   &c,
		Engine: engine,
	}
	return s, nil
}
