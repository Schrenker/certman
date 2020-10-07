package jsonparse

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Settings is a struct representing all settings options available.
type Settings struct {
	Mail             string `json:"mail"`             //default nil
	ConcurrencyLimit int    `json:"concurrencyLimit"` //default 10
}

//InitSettingsJSON is a function that parses settings.json file to Settings struct
func InitSettingsJSON(path string) (*Settings, error) {
	var settings Settings

	err := parseJSONFile(path, &settings)
	if err != nil {
		return &Settings{
			Mail:             "",
			ConcurrencyLimit: 10,
		}, err
	}

	return &settings, nil
}

//InitHostsJSON parses json hosts file into string:[]string map
func InitHostsJSON(path string) (map[string][]string, error) {
	var hosts map[string][]string

	err := parseJSONFile(path, &hosts)
	if err != nil {
		return nil, err
	}

	return hosts, nil

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
