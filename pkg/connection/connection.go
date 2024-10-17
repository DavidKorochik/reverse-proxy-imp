package connection

import (
	"net"

	"davidk/reverse-proxy-imp/pkg/consts"
	"davidk/reverse-proxy-imp/pkg/errors"
)

type Connection struct {
	connStr string
	conn    net.Conn
}

func NewConnection(connString string) (*Connection, error) {
	conn, err := net.Dial(consts.TCPConnection, connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize new TCP connection")
	}

	return &Connection{
		connStr: connString,
		conn:    conn,
	}, nil
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
