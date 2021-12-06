package sc

import (
	"encoding/json"
	"fmt"
	"math"

	g "github.com/Yessir-kim/sessionController/var"
)

type packet struct {
	ID			int
	Sequence	uint
	Payload		[]byte
}

func (s *sessionManager) Read(buf []byte) (int) {
	return s.buffer.read(buf)
}

func (s *sessionManager) Write(buf []byte) (int, error) {
	start, end, path := 0, 0, 0

	for start < len(buf) {
		if start + g.PAYLOAD_SIZE < len(buf) {
			end = start + g.PAYLOAD_SIZE
			fmt.Printf("(Start, End) : (%d, %d)\n", start, end)
		} else {
			end = len(buf)
			fmt.Printf("(Start, End) : (%d, %d)\n", start, end)
		}

		// s.streamList == # of stream
		path = s.seq % len(s.streamList)

		pkt := packet{
			ID: path,
			Sequence: uint(s.seq),
			Payload: buf[start:end],
		}

		fmt.Printf("Payload size : %d\n", len(pkt.Payload))

		bytes, err := marshal(pkt)
		if err != nil {
			fmt.Printf("Marshaling err : %s\n", err)
			panic(err)
		}

		fmt.Printf("%s\n", string(bytes))

		_, err = s.streamList[path].Write(bytes)
		if err != nil {
			fmt.Printf("Writing err : %s\n", err)
			panic(err)
		}
		fmt.Printf("Stream[%d] bytes size : %d\n", path, len(bytes))

		start = end
		s.seq++
	}

	return len(buf), nil

}

func marshal(pkt packet) ([]byte, error) {

	/* prerequisite */
	// First letter in the field name must be capitalized
	// ID, Squence, Payload

	b, err := json.Marshal(pkt)
	if err != nil {
		return b, err
	}

	return b, nil
}

func unmarshal(data []byte) (packet, error) {
	var pkt packet

	json.Unmarshal(data, &pkt)

	return pkt, nil
}

func estPacketSize() (int) {

	digits := math.Ceil(float64(g.PAYLOAD_SIZE / 3))

	return 4 * int(digits) + 37 + 100 // + free space 100 
}
