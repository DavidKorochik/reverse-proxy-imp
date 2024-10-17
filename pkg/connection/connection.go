package connection

import (
	"net"

	"davidk/reverse-proxy-imp/pkg/errors"
)

type Connection struct {
	connStr string

	conn net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		connStr: conn.RemoteAddr().String(),
		conn:    conn,
	}
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

func (c *Connection) Read(bytes []byte) (int, error) {
	return c.conn.Read(bytes)
}

func (c *Connection) Addr() string {
	return c.connStr
}

func GetConnectionsBatch() ([]net.Conn, error) {
	conn1, err := net.Dial("tcp", "0:0:0:0:8000")
	if err != nil {
		return nil, errors.Wrap(err, "from conn1")
	}

	conn2, err := net.Dial("tcp", "0:0:0:0:8001")
	if err != nil {
		return nil, errors.Wrap(err, "from conn2")
	}

	conn3, err := net.Dial("tcp", "0:0:0:0:8002")
	if err != nil {
		return nil, errors.Wrap(err, "from conn3")
	}

	return []net.Conn{conn1, conn2, conn3}, nil
}
