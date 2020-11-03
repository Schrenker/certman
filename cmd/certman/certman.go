package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/schrenker/certman/internal/certutils"
	"github.com/schrenker/certman/internal/jsonparse"
	"github.com/schrenker/certman/internal/queue"
	"github.com/schrenker/certman/internal/sendmail"
)

func main() {
	hostsFlag := flag.String("h", "./configs/hosts.json", "Specify location for hosts json file. Defaults to ./config/hosts")
	settingsFlag := flag.String("s", "./config/settings.json", "Specify location for settings json file. Defaults to ./config/settings")
	flag.Parse()

	hosts, err := jsonparse.InitHostsJSON(*hostsFlag)
	if err != nil {
		fmt.Println("No hosts file provided, exiting...")
		os.Exit(20) //code 20 mean no hosts file provided, which is required for program to run
	}

	settings, _ := jsonparse.InitSettingsJSON(*settingsFlag)

	controlGroup := &queue.ControlGroup{}
	queue.EnqueueHosts(hosts, settings.ConcurrencyLimit, controlGroup)
	controlGroup.Wg.Wait()

	finalList := certutils.GetFinishedCertificateList(settings.Days, hosts)

	if len(settings.EmailServer) > 0 {
		sendmail.Sendmail(settings, finalList)
	} else {
		body := sendmail.PrepareBody(settings.Days, finalList)
		fmt.Print(string(body))
	}
}
