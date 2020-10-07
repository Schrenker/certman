package main

//Certificate monitoring utility

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/schrenker/certman/internal/jsonparse"
)

const usage = `Usage...`

func getCertificateExpiryDate(hostname, domain, port string) time.Time {
	conn, err := net.Dial("tcp", hostname+":"+port)
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}

	client := tls.Client(conn, &tls.Config{
		ServerName: domain,
	})
	defer client.Close()

	if err := client.Handshake(); err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	return client.ConnectionState().PeerCertificates[0].NotAfter
}

func main() {

	jsonFile := jsonparse.LoadJSON("./configs/settings.json")
	jsonparse.ParseHostsJSON(jsonFile)
	defer jsonFile.Close()
}

// func readHosts() //Read hosts from JSON file
// func verifyHostname() //Check if hostname is valid and points to an actual server
// func verifyDomain() //Check if domain give is FQDN
// func enqueueHosts() //Put hosts into queue
// func processQueue() //launch goroutines with upper limit
// func verifyCert() //Check if certificate is valid and if it matches domain
// func parseResult() //Parse resulting date
// func evaluateResult() //evaluate if there is risk of certificate expiring soon
// func sendMail() //Send mail if there are risky certificates
