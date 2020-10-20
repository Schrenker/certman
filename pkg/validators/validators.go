package validators

import (
	"regexp"
	"strconv"
	"strings"
)

//CheckIfFQDN check if supplied domain name is a Fully Qualified Domain Name
func CheckIfFQDN(domain string) bool {
	stdDomain := strings.TrimSpace(strings.ToLower(domain))

	if len(stdDomain) == 0 {
		return false
	}

	if stdDomain[len(stdDomain)-1] == '.' {
		stdDomain = stdDomain[:len(stdDomain)-1]
	}

	regex := "^([a-z0-9]+(-[a-z0-9]+)*\\.)+[a-z]{2,}$"

	match, _ := regexp.MatchString(regex, stdDomain)

	return match
}

func CheckIfIPv4(ip string) bool {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return false
	}
	for i := range octets {
		if len(octets[i]) == 0 {
			return false
		}
		if octet, _ := strconv.Atoi(octets[i]); 0 > octet || octet > 255 {
			return false
		}
	}
	return true
}

// func CheckIfValidPort(port interface{}) bool {}
// func CheckIfValidEmail(email string) bool {}
