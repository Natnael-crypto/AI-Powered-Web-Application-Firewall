package utils

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type IPRange struct {
	startIP     uint32
	endIP       uint32
	countryName string
}

var ipRanges []IPRange

func LoadIPRanges(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		startIP := record[0]
		endIP := record[1]
		countryName := record[2]

		startIPInt, err := ipToUint32(startIP)
		if err != nil {
			return fmt.Errorf("invalid start IP %s: %v", startIP, err)
		}
		endIPInt, err := ipToUint32(endIP)
		if err != nil {
			return fmt.Errorf("invalid end IP %s: %v", endIP, err)
		}

		ipRanges = append(ipRanges, IPRange{
			startIP:     startIPInt,
			endIP:       endIPInt,
			countryName: countryName,
		})
	}

	return nil
}

func ipToUint32(ip string) (uint32, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0, fmt.Errorf("invalid IP address")
	}

	ip = parsedIP.To4().String()

	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid IP address")
	}

	var result uint32
	for i, part := range parts {
		partInt, err := strconv.Atoi(part)
		if err != nil {
			return 0, fmt.Errorf("invalid IP address part %s", part)
		}
		result |= uint32(partInt) << (24 - 8*i)
	}

	return result, nil
}

func GetCountryName(ip string) string {
	ipInt, err := ipToUint32(ip)
	if err != nil {
		return "Invalid IP"
	}

	for _, ipRange := range ipRanges {
		if ipInt >= ipRange.startIP && ipInt <= ipRange.endIP {
			return ipRange.countryName
		}
	}

	return "Unknown"
}

