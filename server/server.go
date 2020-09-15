package server

import (
	"fmt"

	"github.com/mirror-media/yt-relay/config"
	"github.com/mirror-media/yt-relay/relay"
	"github.com/mirror-media/yt-relay/server/route"

	"github.com/gin-gonic/gin"
)

type Server struct {
	conf   config.Conf
	engine *gin.Engine
}

func (s *Server) Run() error {
	return s.engine.Run(fmt.Sprintf("%s:%d", s.conf.Address, s.conf.Port))
}

func New(c config.Conf) (*Server, error) {

	engine := gin.Default()

	relayService, err := relay.New(c.ApiKey)
	if err != nil {
		return nil, err
	}

	_ = route.Set(engine, relayService)

	s := &Server{
		conf:   c,
		engine: engine,
	}
	return s, nil
}
