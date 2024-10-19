package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"davidk/reverse-proxy-imp/pkg/connection"
	"davidk/reverse-proxy-imp/pkg/consts"
	"davidk/reverse-proxy-imp/pkg/errors"
	"davidk/reverse-proxy-imp/pkg/lb"
	"davidk/reverse-proxy-imp/pkg/packets"
)

type Stats struct {
	activeConnections int
	totalConnections  int
}

type Proxy struct {
	port int

	ctx      context.Context
	stats    Stats
	servers  []net.Conn
	ready    chan net.Conn
	listener net.Listener
}

func NewProxy(servers []net.Conn, port int) *Proxy {
	return &Proxy{
		port:    port,
		ctx:     context.Background(),
		servers: servers,
		ready:   make(chan net.Conn, 1),
	}
}

func (p *Proxy) Run() error {
	listener, err := net.Listen(consts.TCPConnection, fmt.Sprintf("%s:%d", consts.LocalHost, p.port))
	if err != nil {
		return errors.WrapF(err, "failed listening on local host and port %s", p.port)
	}

	p.listener = listener
	go p.listen()

	for {
		select {
		case conn := <-p.ready:
			if err = p.Add(connection.NewConnection(conn)); err != nil {
				return err
			}
		case <-p.ctx.Done():
			p.stats.activeConnections--
			break
		}
	}
}

func (p *Proxy) Add(conn *connection.Connection) error {
	p.stats.activeConnections++
	p.stats.totalConnections++

	var data []byte
	if _, err := conn.Read(data); err != nil {
		p.stats.activeConnections--
		return err
	}

	if err := p.deliverPacketToServer(packets.NewPacket(conn, data)); err != nil {
	}

	return nil
}

func (p *Proxy) listen() error {
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			return errors.Wrap(err, "accepting TCP connections")
		}

		p.ready <- conn
	}
}

func (p *Proxy) deliverPacketToServer(packet *packets.Packet) error {
	loadBalancer := lb.NewLoadBalance(p.servers)

	req := &http.Request{}
	if addr := loadBalancer.Next(); addr != "" && loadBalancer.IsHealthy() {
		newReq, err := p.getModifiedRequest(packet)
		if err != nil {
			p.stats.activeConnections--
			return err
		}

		req = newReq
	}

	if err := loadBalancer.ServeHTTP(req); err != nil {
		p.stats.activeConnections--
		return err
	}
}

func (p *Proxy) getModifiedRequest(packet *packets.Packet) (*http.Request, error) {
	var newReq *http.Request
	if err := json.Unmarshal(packet.Data(), &newReq); err != nil {
		p.stats.activeConnections--
		return nil, err
	}

	return newReq, nil
}
