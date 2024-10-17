package packets

type Packet struct {
	data []byte
	len  int64
}

func NewPacket(data []byte) *Packet {
	return &Packet{
		data: data,
		len:  int64(len(data)),
	}
}

func (p *Packet) Len() int64 {
	return p.len
}

// Data - only the packet data reaches the end destination
func (p *Packet) Data() []byte {
	return p.data
}
