package sc

import (
	"encoding/json"
	"fmt"

	g "github.com/Yessir-kim/sessionController/var"
)

type packet struct {
	ID       int
	Total	 int
	Sequence uint
	Payload  []byte
}

func (s *sessionManager) Read(buf []byte) int {
	return s.buffer.read(buf)
}

func (s *sessionManager) Write(buf []byte, w1 int, w2 int) (int, error) {
	start, end, path := 0, 0, 0

	// weight of streams
	stream1 := len(buf) / (w1 + w2) * w1
	stream2 := len(buf) - stream1 // remainning data 

	for start < len(buf) {
		if start + g.PAYLOAD_SIZE < len(buf) {
			end = start + g.PAYLOAD_SIZE
			fmt.Printf("(Start, End) : (%d, %d)\n", start, end)
		} else {
			end = len(buf)
			fmt.Printf("(Start, End) : (%d, %d)\n", start, end)
		}

		// s.streamList == # of stream
		if len(s.streamList) == 1 { // single path connection
			path = 0
		} else {
			path = s.seq % len(s.streamList) // select the stream based on packet sequence

			if w1 != w2 {
				if path == 0 {
					if stream1 < 0 {
						path = 1
						stream2 -= (end - start)
					}
					stream1 -= (end - start)
				} else {
					if stream2 < 0 {
						path = 0
						stream1 -= (end - start)
					}
					stream2 -= (end - start)
				}
			}
		}

		pkt := packet{
			ID:       path,
			Total:	  len(buf),
			Sequence: uint(s.seq),
			Payload:  buf[start:end],
		}

		fmt.Printf("Payload size : %d (Write)\n", len(pkt.Payload))

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
		fmt.Printf("Stream[%d] bytes size : %d (Write)\n", path, len(bytes))

		start = end
		s.seq++
	}

	return len(buf), nil
}

func marshal(pkt packet) ([]byte, error) {

	/* prerequisite */
	// First letter in the field name must be capitalized
	// ID, Total, Squence, Payload

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
