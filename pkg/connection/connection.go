package connection

import (
	"net"
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
