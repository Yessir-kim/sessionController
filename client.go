package sc

import(
		"crypto/tls"

		client "github.com/lucas-clemente/quic-go"
)

func DialAddr(
	addr string,
	tlsConf *tls.Config,
	config *client.Config,
) (client.Session, error) {
		return client.DialAddr(addr, tlsConf, config)
}
