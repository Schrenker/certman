package main

//Certificate monitoring utility

import (
	"os"

	"github.com/schrenker/certman/internal/jsonparse"
	"github.com/schrenker/certman/internal/queue"
)

const usage = `Usage...`

func main() {
	hosts, err := jsonparse.InitHostsJSON("./configs/hosts.json")
	if err != nil {
		os.Exit(20) //code 20 mean no hosts file provided, which is required for program to run
	}

	settings, _ := jsonparse.InitSettingsJSON("./configs/settings.json")

	queue.EnqueueHosts(hosts, settings)
}

// func verifyHostname() //Check if hostname is valid and points to an actual server
// func verifyDomain() //Check if domain give is FQDN
// func enqueueHosts() //Put hosts into queue
// func processQueue() //launch goroutines with upper limit
// func verifyCert() //Check if certificate is valid and if it matches domain
// func parseResult() //Parse resulting date
// func evaluateResult() //evaluate if there is risk of certificate expiring soon
// func sendMail() //Send mail if there are risky certificates
