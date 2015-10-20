package server

import (
	"../payload"
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func is_client_finished(bytes []byte, complete bool) bool {
	// if bytes is of length two with contents {0x00, 0xFF}
	// then that is the close connection message
	// from client that it is finished so
	// return from function to start
	// listening for another client
	if complete && len(bytes) == 2 && bytes[0] == 0x00 {
		return true
	}
	return false
}

func handle_bytes(
	p *payload.Payload,
	bytes []byte,
	complete bool,
) (*payload.Payload, bool) {
	//var p *payload.Payload
	// if last bytes recieved finished the message
	// intialize a new payload object
	if complete {
		p = payload.InitPayload(bytes)
	}
	// copy bytes and update if we are finished
	// for this object and wait for a new
	// payload, or for the exit signal to close connection

	return p, p.AddBytes(bytes)
}

// main worker thread that handles communication with the client
// bytes are parsed into separate payloads and passed
// to a consumer thread
func Communicate(conn net.Conn, payload_channel chan *payload.Payload) {
	fmt.Println("\nopen connection", time.Now(), "\n")
	defer conn.Close()

	var p *payload.Payload
	complete := true

	// create buffered io reader object with connection
	buffer := bufio.NewReader(conn)
	for {
		// read bytes until 0xFF is encountered
		// the first byte sent will be the length
		// so all bytes will be read into correct payload object
		bytes, error := buffer.ReadBytes(0xFF)
		if error == nil && len(bytes) > 0 {
			//fmt.Println(bytes)
			if is_client_finished(bytes, complete) {
				fmt.Println("\nclose connection", time.Now(), "\n")
				return
			}
			p, complete = handle_bytes(p, bytes, complete)
			if complete {
				// send completed payload across channel to consumer
				payload_channel <- p
			}
		}
	}
}

func Triangulator(channel chan *payload.EddyStoneUID) {

	for {
		e := <-channel
		instance := binary.BigEndian.Uint32(e.Instance[2:len(e.Instance)])
		fmt.Print(e.Rssi, instance, "\n")
	}
}

func Payload_consumer(
	input_channel chan *payload.Payload,
	output_channel chan *payload.EddyStoneUID,
	uid []byte,
) {
	for {
		p := <-input_channel
		adv := payload.InitAdvertisement(p.Data)
		valid, frame := payload.ParseEddyStone(adv)
		if valid {
			switch frame.(type) {
			case *payload.EddyStoneUID:
				output_channel <- frame.(*payload.EddyStoneUID)
			default:
			}
		}
	}
}
