package payload

import (
	"fmt"
)

// represents a single complete message
type Payload struct {
	id     int
	Length int
	Data   []byte
}

func (p *Payload) AddBytes(bytes []byte) bool {
	if p.Length == len(p.Data)+len(bytes) {
		// after input bytes are appended payload
		// will be complete
		p.Data = append(p.Data, bytes...)
		p.Data = p.Data[1 : len(p.Data)-1]
		return true
	} else {
		// more data will be recieved
		p.Data = append(p.Data, bytes...)
		return false
	}
}

func (p *Payload) Print() {
	//fmt.Println("rssi", p.Rssi)
	fmt.Println("size", p.Length)
	fmt.Println("data", p.Data, "\n")
}

func InitPayload(bytes []byte) *Payload {
	// create new Payload object with correct defaults
	return &Payload{Length: int(bytes[0]), Data: []byte{}}
}