package sc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	client "github.com/lucas-clemente/quic-go"
	//g "github.com/Yessir-kim/sessionController/var"
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

	fmt.Printf("Remote Address : %s\n", addr)

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

			fmt.Printf("\t\tFor %s address...\n", each.nicIP)

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

			fmt.Printf("\t\tSession Creation!\n")

			stream, err := session.OpenStreamSync(context.Background())
			if err != nil {
				return sessionManager{}, err
			}
			s.streamList = append(s.streamList, stream)

			fmt.Printf("\t\tStream Creation! [%T]\n", stream)

			go func() {
				for {
					// unknown packet size
					buf := make([]byte, 1410)

					n, err := stream.Read(buf)
					if err != nil {
						fmt.Printf("Client Read() error : %s\n", err)
					}

					fmt.Printf("Client Read() data size : %d\n", n)
				}
			}()
		}
	}

	return s, nil
}

func getNICInfo() (nicInfoList, error) {

	nic1 := nicInfo{"wifi", "192.168.0.66"}
	nic2 := nicInfo{"ethernet", "203.252.121.211"}

	n := nicInfoList{[]nicInfo{nic1, nic2}}

	return n, nil
}
