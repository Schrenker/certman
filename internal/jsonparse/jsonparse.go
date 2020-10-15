package jsonparse

import (
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

//Settings is a struct representing all settings options available.
type Settings struct {
	EmailAddr        string   `json:"emailAddr"`        //default ""
	EmailPass        string   `json:"emailPass"`        //default ""
	EmailServer      string   `json:"emailServer"`      //default ""
	EmailPort        string   `json:"emailPort"`        //default ""
	ConcurrencyLimit uint     `json:"concurrencyLimit"` //default 20
	Days             []uint16 `json:"days"`             //default [1, 7, 14]
}

//Vhost ...
type Vhost struct {
	Hostname    string
	Domain      string
	Port        string
	Certificate *x509.Certificate
	Error       error
}

//InitSettingsJSON is a function that parses settings.json file to Settings struct
func InitSettingsJSON(path string) (*Settings, error) {
	var settings Settings

	err := parseJSONFile(path, &settings)
	if err != nil {
		return &Settings{
			EmailAddr:        "",
			EmailPass:        "",
			EmailServer:      "",
			EmailPort:        "",
			ConcurrencyLimit: 20,
			Days:             []uint16{1, 7, 14},
		}, err
	}

	return &settings, nil
}

//InitHostsJSON parses json hosts file into string:[]string map
func InitHostsJSON(path string) ([]*Vhost, error) {
	var hosts map[string][]string

	err := parseJSONFile(path, &hosts)
	if err != nil {
		return nil, err
	}

	return hostsMapToStructArray(hosts), nil
}

//ParseJSON
func parseJSONFile(path string, v interface{}) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	parsedJSONFile, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(parsedJSONFile, &v)
	if err != nil {
		return err
	}

	return nil
}

func hostsMapToStructArray(hosts map[string][]string) []*Vhost {
	vhostArray := make([]*Vhost, amountOfDomainsChecked(hosts))

	i := 0
	for hostname, domains := range hosts {

		for j := range domains {
			if len(domains[j]) == 0 {
				continue
			}

			domain, port := splitPortAndDomain(domains[j])
			vhostArray[i] = &Vhost{
				Hostname: hostname,
				Domain:   domain,
				Port:     port,
			}
			i++
		}
	}
	return vhostArray
}

func amountOfDomainsChecked(hosts map[string][]string) int {
	cumulatedLength := 0
	for _, v := range hosts {
		cumulatedLength += len(v)
	}
	return cumulatedLength
}

func splitPortAndDomain(domain string) (string, string) {
	split := strings.Split(domain, ":")
	if len(split) == 1 {
		return split[0], "443"
	}
	return split[0], split[1]
}
