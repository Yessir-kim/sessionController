package temp

import (
		"crypto/tls"

		server "github.com/lucas-clemente/quic-go"
)

func ListenAddr(addr string, tlsConf *tls.Config, config *server.Config) (server.Listener, error) {
	return server.ListenAddr(addr, tlsConf, config)
}
