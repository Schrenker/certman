package queue

import (
	"strings"
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

	for _, v := range domains {
		domain, port := splitDomainAndPort(v)
		certutils.VerifyCertificates(hostname, domain, port)
	}

	<-limit
}

func splitDomainAndPort(domain string) (string, string) {
	split := strings.Split(domain, ":")
	if len(split) < 2 {
		return split[0], "443"
	}
	return split[0], split[1]
}
