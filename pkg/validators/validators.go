package validators

import (
	"fmt"
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

//CheckIfIPv4 validates if correct ipv4 is provided
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

//CheckIfValidPort validates if given port is correct. Currently support string and int format
func CheckIfValidPort(port interface{}) bool {
	switch v := port.(type) {
	case int:
		if v < 0 && v <= 65535 {
			return true
		}
		return false
	case string:
		numeric, _ := strconv.Atoi(v)
		if numeric < 0 && numeric <= 65535 {
			return true
		}
		return false
	default:
		fmt.Println("Unknown type")
		return false
	}
}
