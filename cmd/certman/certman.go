package main

import (
	"os"

	"github.com/schrenker/certman/internal/jsonparse"
	"github.com/schrenker/certman/internal/queue"
)

func main() {
	hosts, err := jsonparse.InitHostsJSON("./configs/hosts.json")
	if err != nil {
		os.Exit(20) //code 20 mean no hosts file provided, which is required for program to run
	}

	settings, _ := jsonparse.InitSettingsJSON("./configs/settings.json")

	controlGroup := &queue.ControlGroup{}

	queue.EnqueueHosts(hosts, settings, controlGroup)
	controlGroup.Wg.Wait()
	// errors := certutils.GetInvalidCertificatesSlice(hosts)
	// for i := range errors {
	// 	fmt.Println(errors[i].Error)
	// }
	// certsThirty := certutils.GetCertsExpiringInDays(30, hosts)
	// for i := range certsThirty {
	// 	fmt.Printf("%v:%v:%v - %v - %v\n", certsThirty[i].Hostname, certsThirty[i].Domain, certsThirty[i].Port, certsThirty[i].Certificate.NotAfter, certsThirty[i].Error)
	// }
}

// func verifyCert() //Check if certificate is valid and if it matches domain
// func parseResult() //Parse resulting date
// func evaluateResult() //evaluate if there is risk of certificate expiring soon
// func sendMail() //Send mail if there are risky certificates
