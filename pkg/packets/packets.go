package packets

import (
	"davidk/reverse-proxy-imp/pkg/connection"
)

type Packet struct {
	data []byte
	len  int64

	associatedConn connection.Connection
}

func NewPacket(associatedConn *connection.Connection, data []byte) *Packet {
	return &Packet{
		data:           data,
		len:            int64(len(data)),
		associatedConn: *associatedConn,
	}
}

func (p *Packet) Len() int64 {
	return p.len
}

// Data - only the packet data reaches the end destination
func (p *Packet) Data() []byte {
	return p.data
}
