package main

import (
	"os"

	"github.com/schrenker/certman/internal/certutils"
	"github.com/schrenker/certman/internal/jsonparse"
	"github.com/schrenker/certman/internal/queue"
	"github.com/schrenker/certman/internal/sendmail"
)

func main() {
	hosts, err := jsonparse.InitHostsJSON("./configs/hosts.json")
	if err != nil {
		os.Exit(20) //code 20 mean no hosts file provided, which is required for program to run
	}

	settings, _ := jsonparse.InitSettingsJSON("./configs/settings.json")

	controlGroup := &queue.ControlGroup{}
	queue.EnqueueHosts(hosts, settings.ConcurrencyLimit, controlGroup)
	controlGroup.Wg.Wait()

	finalList := certutils.GetFinishedCertificateList(settings.Days, hosts)

	sendmail.Sendmail(settings, finalList)
}

// func verifyCert() //Check if certificate is valid and if it matches domain
// func sendMail() //Send mail if there are risky certificates
