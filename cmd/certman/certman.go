package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/schrenker/certman/internal/certutils"
	"github.com/schrenker/certman/internal/jsonparse"
	"github.com/schrenker/certman/internal/output"
	"github.com/schrenker/certman/internal/queue"
)

func main() {
	hostsFlag := flag.String("h", "./configs/hosts.json", "Specify location for hosts json file. Defaults to ./config/hosts")
	settingsFlag := flag.String("s", "./configs/settings.json", "Specify location for settings json file. Defaults to ./config/settings")
	mailFlag := flag.Bool("m", false, "Send mail with credentials used in settings.json file. Default to false")
	flag.Parse()

	hosts, err := jsonparse.InitHostsJSON(*hostsFlag)
	if err != nil {
		fmt.Println("No hosts file provided, exiting...")
		os.Exit(20)
	}

	settings, _ := jsonparse.InitSettingsJSON(*settingsFlag)

	controlGroup := &queue.ControlGroup{}
	queue.EnqueueHosts(hosts, settings.ConcurrencyLimit, controlGroup)
	controlGroup.Wg.Wait()

	finalList := certutils.GetFinishedCertificateList(settings.Days, hosts)

	if *mailFlag {
		err := output.Sendmail(settings, finalList)
		if err != nil {
			fmt.Println(err)
			os.Exit(30)
		}
	} else {
		body := output.PrepareBody(settings.Days, finalList)
		fmt.Print(body)
	}
}
