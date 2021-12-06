package sc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	//g "github.com/Yessir-kim/sessionController/var"
	server "github.com/lucas-clemente/quic-go"
)

type sessionManager struct {
	// server listener
	listener    server.Listener
	address     string
	sessionList []server.Session
	streamList  []server.Stream
	buffer		*rebuffer
	seq			int
}

func ListenAddr(addr string, tlsConf *tls.Config, config *server.Config) sessionManager {
	lis, err := server.ListenAddr(addr, tlsConf, config)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("Listener Creation\n")

	sessManager := sessionManager{
		listener:    lis,
		address:     addr,
		sessionList: make([]server.Session, 0),
		streamList:  make([]server.Stream, 0),
		buffer: New(),
		seq: 0,
	}

	return sessManager
}

func (s *sessionManager) Accept(ctx context.Context) {
	go s.accept(ctx)
	// blocking until one session is created
	for len(s.sessionList) == 0 {}
}

func (s *sessionManager) accept(ctx context.Context) {
	for {
		sess, err := s.listener.Accept(ctx)
		if err != nil {
			panic(err)
		}
		s.sessionList = append(s.sessionList, sess)

		fmt.Printf("Session Creation!\n")
		//fmt.Printf("%s\n", sess.LocalAddr())
		//fmt.Printf("%s\n", sess.RemoteAddr())

		go func() {
			stream, err := sess.AcceptStream(ctx)
			if err != nil {
				panic(err)
			}
			s.streamList = append(s.streamList, stream)

			fmt.Printf("Stream Creation!\n")

			dec := json.NewDecoder(stream)
			var p packet

			for {
				// buf := make([]byte, estPacketSize())

				if err := dec.Decode(&p); err != nil {
					continue
				} else {
					fmt.Printf("\tPacket id : %d\n", p.ID)
					fmt.Printf("\tPacket seq : %d\n", p.Sequence)
					fmt.Printf("\tPacket payload size : %d\n", len(p.Payload))
				}

				/*
				n, err := stream.Read(buf)
				if err != nil {
					fmt.Printf("Server Read() error : %s\n", err)
				}

				pkt, err := unmarshal(buf[:n])
				
				fmt.Printf("Read() %d size data\n", n)
				fmt.Printf("\tPacket id : %d\n", pkt.ID)
				fmt.Printf("\tPacket seq : %d\n", pkt.Sequence)
				fmt.Printf("\tPakcet payload size : %d\n", len(pkt.Payload))
				*/
				for !s.buffer.store(p.Payload[:len(p.Payload)], int(p.Sequence)) {}
			}
		}()
	}
}
