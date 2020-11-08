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
//
//First argument is a server to which we are connecting
//
//Second argument is domain of which we are checking certificate
//
//Third argument is port to which we are connecting
func getCertificate(hostname, domain, port string) (*x509.Certificate, error) {
	conn, err := net.Dial("tcp", hostname+":"+port)
	if err != nil {
		return &x509.Certificate{}, err
	}
	defer conn.Close()

	client := tls.Client(conn, &tls.Config{
		ServerName: domain,
	})
	defer client.Close()

	if err := client.Handshake(); err != nil {
		return &x509.Certificate{}, err
	}

	return client.ConnectionState().PeerCertificates[0], nil
}

//GetInvalidCertificatesSlice returns a slice containing only invalid certificates and entries which failed
//
//Function takes one argument, which is a list of all vhosts
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

//GetFinishedCertificateList composes array of vhost arrays, containing all fetches certificates, categorized via "days" breakpoints.
//
//First argument is days array, which break up list into manageable chunks
//
//Second argument is array of all vhosts
func GetFinishedCertificateList(days []uint16, vhosts []*jsonparse.Vhost) [][]*jsonparse.Vhost {
	finalList := make([][]*jsonparse.Vhost, 0)

	for i := range days {
		finalList = append(finalList, getCertsExpiringInDays(days[i], vhosts))
	}

	cutDuplicates(finalList)

	finalList = append(finalList, getInvalidCertificatesSlice(vhosts))

	return finalList
}
