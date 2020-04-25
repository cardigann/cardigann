package main

import (
	"encoding/json"
	"fmt"
	_ "log"
	"os"

	"github.com/cardigann/cardigann/server"
	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
)

type programOpts struct {
	UserService bool
}

type program struct {
	exit    chan struct{}
	service service.Service
	logger  service.Logger
}

func newProgram(opts programOpts) (*program, error) {
	svcConfig := &service.Config{
		Name:        "Cardigann",
		DisplayName: "Cardigann Proxy",
		Description: "Cardigann Torrent Indexer Proxy",
		Option: service.KeyValue{
			"RunAtLoad":   true,
			"UserService": opts.UserService,
		},
		Arguments: []string{
			"service", "run",
		},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, err
	}

	prg.service = s
	return prg, nil
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		p.logger.Info("Running interactively")
	}
	p.exit = make(chan struct{})

	go p.run()
	return nil
}

func (p *program) run() error {
	p.logger.Infof("Running service via %v.", service.Platform())

	conf, err := newConfig()
	if err != nil {
		return err
	}

	s, err := server.New(conf, Version)
	if err != nil {
		return err
	}

	go s.Listen()

	// block until exit
	<-p.exit
	return nil
}

func (p *program) Stop(s service.Service) error {
	p.logger.Info("Shutting down cardigann")
	p.exit <- struct{}{}
	return nil
}

type serviceLogHook struct {
	service.Logger
}

func (hook *serviceLogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		hook.Logger.Error(line)
	case logrus.FatalLevel:
		hook.Logger.Error(line)
	case logrus.ErrorLevel:
		hook.Logger.Error(line)
	case logrus.WarnLevel:
		hook.Logger.Warning(line)
	case logrus.InfoLevel:
		hook.Logger.Info(line)
	case logrus.DebugLevel:
		hook.Logger.Info(line)
	}

	return nil
}

func (hook *serviceLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type serviceLogFormatter struct {
}

func (f *serviceLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return []byte(fmt.Sprintf("%s %s", entry.Message, serialized)), nil
}
