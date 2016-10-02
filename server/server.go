package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/logger"
)

// Server is an http server which wraps the Handler
type Server struct {
	Bind, Port, Passphrase string
	version                string
	config                 config.Config
}

func New(conf config.Config, version string) (*Server, error) {
	bind, err := config.GetGlobalConfig("bind", "0.0.0.0", conf)
	if err != nil {
		return nil, err
	}

	port, err := config.GetGlobalConfig("port", "5060", conf)
	if err != nil {
		return nil, err
	}

	passphrase, err := config.GetGlobalConfig("passphrase", "", conf)
	if err != nil {
		return nil, err
	}

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if version == "" {
		version = "dev"
	}

	return &Server{
		Bind:       bind,
		Port:       port,
		Passphrase: passphrase,
		config:     conf,
		version:    version,
	}, nil
}

func (s *Server) Listen() error {
	logger.Logger.Infof("Cardigann %s", s.version)

	listenOn := fmt.Sprintf("%s:%s", s.Bind, s.Port)
	logger.Logger.Infof("Listening on %s", listenOn)

	h, err := NewHandler(Params{
		Passphrase: s.Passphrase,
		Config:     s.config,
		Version:    s.version,
	})
	if err != nil {
		return err
	}

	return http.ListenAndServe(listenOn, h)
}
