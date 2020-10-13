package certutils

import (
	"crypto/tls"
	"crypto/x509"
	"net"

	"github.com/schrenker/certman/internal/jsonparse"
)

//VerifyCertificates launches certificate retrieval and test process.
func VerifyCertificates(vhost *jsonparse.Vhost) {
	vhost.Certificate, vhost.Error = getCertificate(vhost.Hostname, vhost.Domain, vhost.Port)
}

//getCertificateExpiryDate connects to a server specified by hostname argument and port, and then creates tcp.Client connection to specified domain.
//It is done so to prevent retrieving unimportant certificate like Cloudflare.
//This method of connection also allows to check vhost's certificate that is otherwise unavailable on the network via domain name.
func getCertificate(hostname, domain, port string) (*x509.Certificate, error) {
	conn, err := net.Dial("tcp", hostname+":"+port)
	defer conn.Close()

	if err != nil {
		return &x509.Certificate{}, err
	}

	client := tls.Client(conn, &tls.Config{
		ServerName: domain,
	})
	defer client.Close()

	if err := client.Handshake(); err != nil {
		return &x509.Certificate{}, err
	}

	return client.ConnectionState().PeerCertificates[0], nil
}
