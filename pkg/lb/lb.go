package lb

import (
	"sync"

	"davidk/reverse-proxy-imp/pkg/proxy"
	"davidk/reverse-proxy-imp/pkg/server"
)

type Service interface{}

type LoadBalancer struct {
	currentIndex int

	servers []*server.Server
	proxy   proxy.Proxy
	mux     sync.RWMutex
}

func NewLoadBalancer(servers []*server.Server) Service {
	return &LoadBalancer{
		servers: servers,
	}
}
