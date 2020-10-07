package main

//Certificate monitoring utility

import (
	"fmt"

	"github.com/schrenker/certman/internal/jsonparse"
)

const usage = `Usage...`

func main() {
	hosts, err := jsonparse.InitHostsJSON("./configs/hosts.json")
	fmt.Println(hosts, err)
	settings, err := jsonparse.InitSettingsJSON("./configs/settings.json")
	fmt.Println(settings, err)
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
