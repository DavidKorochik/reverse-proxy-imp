package server

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"sync"

	"davidk/reverse-proxy-imp/pkg/consts"
)

type Service interface {
	IsHealthy() bool
	IncreaseActiveConnections()
	IncreaseTotalConnections()
	GetServerStats() Stats
	GetURL() *url.URL
	ServeHTTP(rw http.ResponseWriter, req *http.Request)
	PerformServerHealthCheck(ctx context.Context, url *url.URL)
}

type Stats struct {
	TotalConnections  int
	ActiveConnections int
}

type Server struct {
	isHealthy   bool
	serverStats Stats

	mux sync.RWMutex // using read-write mutex since we have specific functions for reading and writing servers stats
	url *url.URL
}

func NewServer(url *url.URL) Service {
	return &Server{
		url: url,
	}
}

func (s *Server) IsHealthy() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.isHealthy
}

func (s *Server) IncreaseActiveConnections() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.serverStats.ActiveConnections++
}

func (s *Server) IncreaseTotalConnections() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.serverStats.TotalConnections++
}

func (s *Server) GetServerStats() Stats {
	return s.serverStats
}

func (s *Server) GetURL() *url.URL {
	return s.url
}

func (s *Server) PerformServerHealthCheck(ctx context.Context, url *url.URL) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, consts.TCPConnection, url.Host)
	if err != nil {
		s.mux.RLock()
		defer s.mux.RUnlock()

		s.isHealthy = false
		return
	}
	defer conn.Close()

	s.isHealthy = true
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// send http request through proxy server
}
