package certutils

import (
	"crypto/tls"
	"net"
	"time"
)

//GetCertificateExpiryDate is a function that makes connection, and probes domains certificates via tls.client
func GetCertificateExpiryDate(hostname, domain, port string) (time.Time, error) {
	conn, err := net.Dial("tcp", hostname+":"+port)
	defer conn.Close()

	if err != nil {
		return time.Time{}, err
	}

	client := tls.Client(conn, &tls.Config{
		ServerName: domain,
	})
	defer client.Close()

	if err := client.Handshake(); err != nil {
		return time.Time{}, err
	}

	return client.ConnectionState().PeerCertificates[0].NotAfter, nil
}
