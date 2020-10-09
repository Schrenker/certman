package certutils

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

//VerifyCertificates launches certificate retrieval and test process.
func VerifyCertificates(hostname, domain, port string) {
	date, _ := getCertificateExpiryDate(hostname, domain, port)
	fmt.Printf("%v:%v:%v - %v\n", hostname, domain, port, date)
}

//getCertificateExpiryDate connects to a server specified by hostname argument and port, and then creates tcp.Client connection to specified domain.
//It is done so to prevent retrieving unimportant certificate like Cloudflare.
//This method of connection also allows to check vhost's certificate that is otherwise unavailable on the network via domain name.
func getCertificateExpiryDate(hostname, domain, port string) (time.Time, error) {
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
