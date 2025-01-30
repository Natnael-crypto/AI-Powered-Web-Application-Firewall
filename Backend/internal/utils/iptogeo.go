package utils

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Struct to store an IP range and its corresponding country name
type IPRange struct {
	startIP     uint32
	endIP       uint32
	countryName string
}

// Global variable to hold the IP ranges and country names
var ipRanges []IPRange

// Function to load IP ranges and country names from a CSV file
func LoadIPRanges(filename string) error {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	// Initialize the CSV reader
	reader := csv.NewReader(file)

	// Read the CSV rows and store them
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Parse the IP range and country name from the record
		startIP := record[0]
		endIP := record[1]
		countryName := record[2]

		// Convert the start and end IPs to uint32 for comparison
		startIPInt, err := ipToUint32(startIP)
		if err != nil {
			return fmt.Errorf("invalid start IP %s: %v", startIP, err)
		}
		endIPInt, err := ipToUint32(endIP)
		if err != nil {
			return fmt.Errorf("invalid end IP %s: %v", endIP, err)
		}

		// Add the range and country to the slice
		ipRanges = append(ipRanges, IPRange{
			startIP:     startIPInt,
			endIP:       endIPInt,
			countryName: countryName,
		})
	}

	return nil
}

// Function to convert an IP address to a uint32 integer
func ipToUint32(ip string) (uint32, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0, fmt.Errorf("invalid IP address")
	}

	// Convert the IP to a 4-byte slice
	ip = parsedIP.To4().String()

	// Split the IP into its octets
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid IP address")
	}

	// Convert the octets to an integer
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

// Function to get the country name for a given IP address
func GetCountryName(ip string) string {
	// Convert the IP to uint32
	ipInt, err := ipToUint32(ip)
	if err != nil {
		return "Invalid IP"
	}

	// Check which range the IP belongs to
	for _, ipRange := range ipRanges {
		if ipInt >= ipRange.startIP && ipInt <= ipRange.endIP {
			return ipRange.countryName
		}
	}

	return "Unknown" // Return "Unknown" if IP doesn't match any range
}

