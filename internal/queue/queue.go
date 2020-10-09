package queue

import (
	"fmt"
	"sync"

	"github.com/schrenker/certman/internal/certutils"
	"github.com/schrenker/certman/internal/jsonparse"
)

//EnqueueHosts launches goroutines that take care of checking certificates on hosts
func EnqueueHosts(hosts map[string][]string, settings *jsonparse.Settings) {
	var wg sync.WaitGroup
	wg.Add(len(hosts))

	limit := make(chan struct{}, settings.ConcurrencyLimit) //limit amount of running jobs

	for k, v := range hosts {
		limit <- struct{}{}
		go launchConnection(k, v, &wg, limit)
	}

	wg.Wait()
	close(limit)
}

func launchConnection(hostname string, domains []string, wg *sync.WaitGroup, limit chan struct{}) {
	defer wg.Done()

	for _, domain := range domains {
		date, _ := certutils.GetCertificateExpiryDate(hostname, domain, "443")
		fmt.Printf("%v:%v:443 - %v\n", hostname, domain, date)
	}

	<-limit
}
