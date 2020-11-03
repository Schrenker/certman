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

const usage = `certman - certificate monitoring utility - github.com/schrenker/certman

Usage: certman -h [FILE] -s [FILE] -m [BOOL]
-h - Path to hosts.json file. Defaults to ./hosts.json. This file is required.
-s - Path to settings.json file. Defaults to ./settings.json. Optional
-m - Flag specifying if mail is to be sent. Defaults to false. Optional
`

func main() {
	hostsFlag := flag.String("h", "./hosts.json", "Specify location for hosts json file. Defaults to ./config/hosts")
	settingsFlag := flag.String("s", "./settings.json", "Specify location for settings json file. Defaults to ./config/settings")
	mailFlag := flag.Bool("m", false, "Send mail with credentials used in settings.json file. Default to false")
	flag.Parse()

	hosts, err := jsonparse.InitHostsJSON(*hostsFlag)
	if err != nil {
		fmt.Println(usage)
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
		fmt.Print(output.PrepareBody(settings.Days, finalList))
	}
}
