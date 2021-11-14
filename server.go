package sc

import (
		"crypto/tls"
		"context"
		"fmt"

		server "github.com/lucas-clemente/quic-go"
)

type sessionManager struct {
		// server listener 
		listener server.Listner
		address string
		sessionList []server.Session
		streamList []server.Stream
}

func ListenAddr(addr string, tlsConf *tls.Config, config *server.Config) (sessionManager, error) {
		sessManager := sessionManager{
			server.ListenAddr(addr, tlsConf, config),
			addr,
			make([]server.Session, 0),
			make([]server.Stream, 0),
		}

		return sessManger
}

func (s *sessionManager) Accept(ctx context.Context) error {
		var i int =  0

		for {
			sess, err := s.listener.Accept(ctx)
			if err != nil {
				return err
			}
			s.sessionList = append(s.sessionList, sess)

			stream, err := sess.AcceptStream(ctx)
			if err != nil {
				return err
			}
			s.streamList = append(s.streamList, stream)

			fmt.Printf("Session & Stream [%d] creation\n", i)
			i++

			/*
			go func() {
				for {
					
				}
			}
			*/
		}
}
