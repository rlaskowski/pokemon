package pokemon

import (
	"context"
	"log"
	"os"
	"os/signal"
)

type Service struct {
	HttpServer *HttpServer
	signal     chan os.Signal
	cancel     context.CancelFunc
	ctx        context.Context
}

func NewService() *Service {
	ctx, cancel := context.WithCancel(context.Background())

	service := &Service{
		HttpServer: NewHttpServer(ctx),
		signal:     make(chan os.Signal, 3),
		cancel:     cancel,
		ctx:        ctx,
	}
	return service
}

//Running all services
func (s *Service) Run() error {
	log.Print("Starting all services")

	if err := s.HttpServer.Start(); err != nil {
		return err
	}

	signal.Notify(s.signal, os.Interrupt)
	<-s.signal

	return s.Stop()
}

//Stopping all running services
func (s *Service) Stop() error {
	log.Print("Stopping all services")

	s.cancel()
	close(s.signal)

	if err := s.HttpServer.Stop(); err != nil {
		log.Fatalf("HTTP Server stopping error: %s", err)
		return err
	}

	return nil
}
