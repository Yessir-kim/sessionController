package sc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"

	client "github.com/lucas-clemente/quic-go"
)

type nicInfoList struct {
	lists []nicInfo
}

type nicInfo struct {
	nicType string
	nicIP   string
}

func Dial(addr string, tlsConf *tls.Config) (sessionManager, error) {
	nics, err := getNICInfo()
	if err != nil {
		return sessionManager{}, err
	}

	fmt.Printf("Remote Address : %s (client)\n", addr)

	s := sessionManager{
		listener:    nil,
		address:     addr,
		sessionList: make([]client.Session, 0),
		streamList:  make([]client.Stream, 0),
		buffer: New(),
	}

	for i, each := range nics.lists {

		fmt.Printf("NIC[%d] connection\n", i)

		if each.nicType == "wifi" || each.nicType == "ethernet" {
			udpAddr, err := net.ResolveUDPAddr("udp", addr)
			if err != nil {
				return sessionManager{}, err
			}

			fmt.Printf("For %s address... (client)\n", each.nicIP)

			ip4 := net.ParseIP(each.nicIP).To4()
			udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip4, Port: 0})
			if err != nil {
				return sessionManager{}, err
			}

			session, err := client.Dial(udpConn, udpAddr, addr, tlsConf, nil)
			if err != nil {
				return sessionManager{}, err
			}
			s.sessionList = append(s.sessionList, session)

			fmt.Printf("\tSession Creation! (client)\n")

			go func() {
					stream, err := session.OpenStreamSync(context.Background())
					if err != nil {
						panic(err)
					}

					s.streamList = append(s.streamList, stream)

					fmt.Printf("\tStream[%d] Creation! (client)\n", len(s.streamList))

					dec := json.NewDecoder(stream)
					var p packet

					for
					{
						if err := dec.Decode(&p); err != nil {
							continue
						} else {
							fmt.Printf("\t\tPacket id : %d (client)\n", p.ID)
							fmt.Printf("\t\tPacket seq : %d (client)\n", p.Sequence)
							fmt.Printf("\t\tPacket payload size : %d (client)\n", len(p.Payload))
						}

						for !s.buffer.store(p.Payload[:len(p.Payload)], int(p.Sequence)) {}
					}
			}()
		}
	}

	return s, nil
}

func getNICInfo() (nicInfoList, error) {

	nic1 := nicInfo{"wifi", "0.0.0.0"}
	nic2 := nicInfo{"ethernet", "0.0.0.0"}

	nics := nicInfoList{[]nicInfo{nic1, nic2}}

	return nics, nil
}
