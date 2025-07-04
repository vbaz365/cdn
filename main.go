package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/vbaz365/cdn/datastructure"
)

func main() {
	data := &datastructure.Data{}
	loadRoutingData(data, "./data/routing-data.txt")

	if len(os.Args) < 2 {
		fmt.Println("Missing input")
		return
	}

	_, ipNet, _ := net.ParseCIDR(os.Args[1])
	pop, scope := data.Route(ipNet)
	fmt.Printf("Pop id: %d, Scope: %d\n", pop, scope)
}

// loadRoutingData loads the routing data from a text file into the datastructure
func loadRoutingData(data *datastructure.Data, filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			log.Printf("Skipping invalid line: %s", line)
			continue
		}

		prefixStr := parts[0]
		popStr := parts[1]

		popNum, err := strconv.ParseUint(popStr, 10, 16)
		if err != nil {
			log.Printf("Invalid pop number in line %s: %v", line, err)
			continue
		}

		_, ipnet, err := net.ParseCIDR(prefixStr)
		if err != nil {
			log.Printf("Invalid prefix in line %s: %v", line, err)
			continue
		}

		// Get prefix length
		prefixLen, _ := ipnet.Mask.Size()

		ip := ipnet.IP.To16()
		if ip == nil {
			log.Printf("Invalid IPv6 in line %s", line)
			continue
		}

		// Get higher and lower half of the IPv6 represented as 2 uint64s
		high := binary.BigEndian.Uint64(ip[:8])
		low := binary.BigEndian.Uint64(ip[8:])

		data.InsertRadixNode(high, low, uint8(prefixLen), uint16(popNum))
	}

	hasErrorOccured := scanner.Err()

	if hasErrorOccured != nil {
		log.Fatal(hasErrorOccured)
	}
}
