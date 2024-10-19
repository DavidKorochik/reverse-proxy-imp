package lb

import (
	"net"
	"sync"

	"davidk/reverse-proxy-imp/pkg/consts"
	"davidk/reverse-proxy-imp/pkg/proxy"
	"davidk/reverse-proxy-imp/pkg/server"
)

type Service interface {
	Rotate() *server.Server
	GetNextValidServer() *server.Server
	GetServersPool() int
}

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

func (lb *LoadBalancer) Rotate() *server.Server {
	lb.mux.Lock()
	lb.currentIndex = (lb.currentIndex + 1) % lb.GetServersPool()
	lb.mux.Unlock()

	return lb.servers[lb.currentIndex]
}

func (lb *LoadBalancer) GetNextValidServer() *server.Server {
	for i := 0; i < lb.GetServersPool(); i++ {
		if rotatedServer := lb.Rotate(); rotatedServer.IsHealthy() {
			return rotatedServer
		}
	}

	return nil
}

func (lb *LoadBalancer) GetServersPool() int {
	return len(lb.servers)
}

func (lb *LoadBalancer) Listen(servers chan server.Service) {
	select {
	case s := <-servers:
		listener, err := net.Listen(consts.TCPConnection, s.GetURL().Host)
		if err != nil {
			break
		}
		listener.Accept()
	}
}

func (lb *LoadBalancer) ServeServerRequestToProxy() {
}
