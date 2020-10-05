package main

//Certificate monitoring utility

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

func getCertificateExpiryDate(hostname, addr, port string) time.Time {
	conn, err := net.Dial("tcp", hostname+":"+port)
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}

	client := tls.Client(conn, &tls.Config{
		ServerName: addr,
	})
	defer client.Close()

	if err := client.Handshake(); err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	return client.ConnectionState().PeerCertificates[0].NotAfter
}

func main() {
	fmt.Println()
}
