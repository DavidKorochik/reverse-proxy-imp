package proxy

import (
	"fmt"
	"net"

	"davidk/reverse-proxy-imp/pkg/errors"
	"davidk/reverse-proxy-imp/pkg/packets"
)

type ProxyStats struct {
	activeConnections int
	totalConnections  int
}

type Proxy struct {
	port int

	servers  []net.Conn
	stats    ProxyStats
	listener net.Listener
	ready    chan net.Conn
}

func NewProxy(servers []net.Conn, port int) *Proxy {
	return &Proxy{
		port:    port,
		servers: servers,
		ready:   make(chan net.Conn, 1),
	}
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

func (p *Proxy) Add(conn net.Conn) error {
	p.stats.activeConnections++

	var data []byte
	if _, err := conn.Read(data); err != nil {
		return err
	}
	packet := packets.NewPacket(data)

	return nil
}

func (p *Proxy) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0:0:0:0:%d", p.port))
	if err != nil {
		return errors.WrapF(err, "failed listening on local host and port %s", p.port)
	}

	p.listener = listener

	go func() error {
		if err = p.listen(); err != nil {
			return err
		}

		return nil
	}()

	for {
		select {
		case conn := <-p.ready:
			if err = p.Add(conn); err != nil {
				return err
			}
		}
	}
}
