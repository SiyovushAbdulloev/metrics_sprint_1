package httpserver

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	App     *gin.Engine
	Address string
}

func New(opts ...Option) *Server {
	server := &Server{
		App:     nil,
		Address: "",
	}

	for _, opt := range opts {
		opt(server)
	}

	app := gin.New()

	server.App = app

	return server
}

func (s *Server) Start() error {
	return s.App.Run(s.Address)
}
