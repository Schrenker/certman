package queue

import (
	"sync"

	"github.com/schrenker/certman/internal/certutils"
	"github.com/schrenker/certman/internal/jsonparse"
)

//EnqueueHosts launches goroutines that take care of checking certificates on hosts
func EnqueueHosts(hosts []*jsonparse.Vhost, settings *jsonparse.Settings) {
	var wg sync.WaitGroup
	wg.Add(len(hosts))

	limit := make(chan struct{}, settings.ConcurrencyLimit) //limit amount of running jobs

	for i := range hosts {
		limit <- struct{}{}
		go launchConnection(hosts[i], &wg, limit)
	}

	wg.Wait()
	close(limit)
}

func launchConnection(vhost *jsonparse.Vhost, wg *sync.WaitGroup, limit chan struct{}) {
	defer wg.Done()

	certutils.VerifyCertificates(vhost)

	<-limit
}
