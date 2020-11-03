package queue

import (
	"crypto/x509"
	"fmt"
	"sync"

	"github.com/schrenker/certman/internal/certutils"
	"github.com/schrenker/certman/internal/jsonparse"
	"github.com/schrenker/certman/pkg/validators"
)

//ControlGroup is a struct containing wait group and channel semaphore
type ControlGroup struct {
	Wg    sync.WaitGroup
	limit chan struct{}
}

//EnqueueHosts launches goroutines that take care of checking certificates on hosts
//
//First argument is an array of vhosts in hostname:domain:port format
//
//Second Argument takes number to which concurrent jobs will be limited
//
//Third arguments takes a pointer to ControlGroup struct
func EnqueueHosts(hosts []*jsonparse.Vhost, concurrencyLimit uint8, controlGroup *ControlGroup) {
	controlGroup.Wg.Add(len(hosts))

	controlGroup.limit = make(chan struct{}, concurrencyLimit) //limit amount of running jobs

	for i := range hosts {
		controlGroup.limit <- struct{}{}
		go launchConnection(hosts[i], &controlGroup.Wg, controlGroup.limit)
	}

	close(controlGroup.limit)
}

func launchConnection(vhost *jsonparse.Vhost, wg *sync.WaitGroup, limit chan struct{}) {
	defer wg.Done()

	if !(validators.CheckIfFQDN(vhost.Hostname) || validators.CheckIfIPv4(vhost.Hostname)) {
		vhost.Error = fmt.Errorf("%v is not a valid domain or IP address", vhost.Hostname)

	} else if !validators.CheckIfFQDN(vhost.Domain) {
		vhost.Error = fmt.Errorf("%v is not a valid domain", vhost.Domain)

	} else if !validators.CheckIfValidPort(vhost.Port) {
		vhost.Error = fmt.Errorf("%v is not a valid port", vhost.Port)

	} else {
		certutils.VerifyCertificates(vhost)
	}

	if vhost.Certificate == nil {
		vhost.Certificate = &x509.Certificate{}
	}
	<-limit
}
