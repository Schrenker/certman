package queue

import (
	"fmt"
	"sync"

	"github.com/schrenker/certman/internal/jsonparse"
)

type concurrencyLock struct {
	mutex                 sync.Mutex
	concurrentConnections int
}

//EnqueueHosts
func EnqueueHost(hosts map[string][]string, settings *jsonparse.Settings) {
	var wg sync.WaitGroup

	for k, v := range hosts {
		wg.Add(1)
		go launchConnection(k, v, &wg)
	}

	wg.Wait()
}

func launchConnection(host string, domains []string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println(host)
	for i := range domains {
		fmt.Println(domains[i])
	}
}
