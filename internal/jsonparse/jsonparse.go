package jsonparse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//LoadJSON is a function that loads json file
func LoadJSON(path string) *os.File {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	return jsonFile
}

//ParseHostsJSON parses json hosts file into string:[]string map
func ParseHostsJSON(jsonFile *os.File) map[string][]string {
	var hosts map[string][]string

	parsedJSONFile, _ := ioutil.ReadAll(jsonFile)

	err := json.Unmarshal(parsedJSONFile, &hosts)
	if err != nil {
		fmt.Println(err)
	}

	jsonFile.Close()
	return hosts

}
