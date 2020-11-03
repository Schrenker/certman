package certutils

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"sort"
	"time"

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

//GetInvalidCertificatesSlice ...
func getInvalidCertificatesSlice(vhosts []*jsonparse.Vhost) []*jsonparse.Vhost {
	errors := make([]*jsonparse.Vhost, 0)
	for i := range vhosts {
		if vhosts[i].Error != nil {
			errors = append(errors, vhosts[i])
		}
	}
	return errors
}

func getCertsExpiringInDays(days uint16, vhosts []*jsonparse.Vhost) []*jsonparse.Vhost {
	expiringCerts := make([]*jsonparse.Vhost, 0)
	expiryDate := time.Now().Add(time.Duration(days*24) * time.Hour)
	for i := range vhosts {
		if vhosts[i].Certificate.NotAfter.Before(expiryDate) && vhosts[i].Error == nil {
			expiringCerts = append(expiringCerts, vhosts[i])
		}
	}
	sort.SliceStable(expiringCerts, func(i, j int) bool {
		return expiringCerts[i].Certificate.NotAfter.Sub(expiringCerts[j].Certificate.NotAfter) < 0
	})
	return expiringCerts
}

func cutDuplicates(vhosts [][]*jsonparse.Vhost) {
	length := len(vhosts)
	if length < 2 {
		return
	}
	combinedLength := 0
	for i := 1; i < length; i++ {
		combinedLength = combinedLength + len(vhosts[i-1])
		vhosts[i] = vhosts[i][combinedLength:]
	}
}

func GetFinishedCertificateList(days []uint16, vhosts []*jsonparse.Vhost) [][]*jsonparse.Vhost {
	finalList := make([][]*jsonparse.Vhost, 0)

	for i := range days {
		finalList = append(finalList, getCertsExpiringInDays(days[i], vhosts))
	}

	cutDuplicates(finalList)

	finalList = append(finalList, getInvalidCertificatesSlice(vhosts))

	return finalList
}
