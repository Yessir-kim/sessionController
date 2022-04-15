package sc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	server "github.com/lucas-clemente/quic-go"
)

type sessionManager struct {
	// server listener
	listener    server.Listener
	address     string
	sessionList []server.Connection
	streamList  []server.Stream
	buffer		*rebuffer
	seq			int
	// mp			bool
}

func ListenAddr(addr string, tlsConf *tls.Config, config *server.Config) sessionManager {
	lis, err := server.ListenAddr(addr, tlsConf, config)
	if err != nil {
		fmt.Printf("server ListenAddr() error : %s\n", err)
	}
	fmt.Printf("Listener Creation! (ListenAddr)\n")

	s := sessionManager{
		listener:    lis,
		address:     addr,
		sessionList: make([]server.Connection, 0),
		streamList:  make([]server.Stream, 0),
		buffer: New(),
		seq: 0,
	}

	return s
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

		fmt.Printf("\tSession[%d] Creation! (ListenAddr)\n", len(s.sessionList))

		go func() {
			stream, err := sess.AcceptStream(ctx)
			if err != nil {
				fmt.Println(err)
			}
			s.streamList = append(s.streamList, stream)

			fmt.Printf("\tStream[%d] Creation! (ListenAddr)\n", len(s.streamList))

			dec := json.NewDecoder(stream)
			var p packet

			for
			{
				if err := dec.Decode(&p); err != nil {
					continue
				} else {
					fmt.Printf("\t\tPacket id : %d (ListenAddr)\n", p.ID)
					fmt.Printf("\t\tPacket total size : %d (ListenAddr)\n", p.Total)
					fmt.Printf("\t\tPacket seq : %d (ListenAddr)\n", p.Sequence)
					fmt.Printf("\t\tPacket payload size : %d (ListenAddr)\n", len(p.Payload))
				}

				// never false 
				for !s.buffer.store(p.Payload[:len(p.Payload)], int(p.Sequence), int(p.Total)) {}
			}
		}()
	}
}
