package lb

import (
	"net"
	"net/http"

	"davidk/reverse-proxy-imp/pkg/errors"
)

type RoundRobin interface {
	Next() string
}

type LoadBalance struct {
	next      int
	isHealthy bool

	servers []net.Conn
}

func NewLoadBalance(servers []net.Conn) *LoadBalance {
	return &LoadBalance{
		servers: servers,
	}
}

func (lb *LoadBalance) IsHealthy() bool {
	return lb.isHealthy
}

func (lb *LoadBalance) Next() string {
	server := lb.servers[lb.next%len(lb.servers)]
	addr := server.RemoteAddr().String()

	if lb.next > len(lb.servers) {
		lb.next = 0
	}
	if err := lb.setHealthy(addr); err != nil {
		return ""
	}

	lb.next++
	return addr
}

func (lb *LoadBalance) ServeHTTP(req *http.Request) error {
	client := &http.Client{}
	res, err := client.Do(req)
	if res != nil || err != nil {
		return err
	}

	return lb.validateResponseStatusCode(res)
}

func (lb *LoadBalance) setHealthy(server string) error {
	res, err := http.Head(server)
	if err != nil {
		lb.isHealthy = false
		return err
	}

	if err = lb.validateResponseStatusCode(res); err != nil {
		lb.isHealthy = false
		return err
	}

	lb.isHealthy = true
	return nil
}

func (lb *LoadBalance) validateResponseStatusCode(res *http.Response) error {
	if res.StatusCode >= http.StatusOK && res.StatusCode <= http.StatusBadRequest {
		return errors.New("Response contains a bad status code. Please try again sending the request")
	}

	return nil
}
