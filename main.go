// geoip: utility for finding ISO 3166 country code using csv input database
// By default, it is configured to use db-ip.com IP to Country Lite database in csv format
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
)

// TODO: if address is a subnet mask according to flag mask - find all range matches

const defaultInputFilename = "data.csv"

const (
	defaultInputIpAddressStartIndex = iota
	defaultInputIpAddressEndIndex
	defaultInputCountryIndex
	defaultSubnetMask = 32
)

var (
	geoDB        = flag.String("geodb", defaultInputFilename, "csv file with geoip data")
	ipStartIndex = flag.Int("ipstart", defaultInputIpAddressStartIndex, "ip address start index in csv file (default 0)")
	ipEndIndex   = flag.Int("ipend", defaultInputIpAddressEndIndex, "ip address end index in csv file")
	countryIndex = flag.Int("country", defaultInputCountryIndex, "country index in csv file")
	showDbSize   = flag.Bool("size", false, "show number of records in database")
	checkDB      = flag.Bool("check", false, "check geoip for completeness")
	ipMask       = flag.Int("mask", defaultSubnetMask, "subnet mask")
	wildcard     int
)

type IPRange struct {
	start int64
	end   int64
	cc    string
}

type GeoIP struct {
	ipToCountryCode []IPRange
}

func (g *GeoIP) Size() int {
	return len(g.ipToCountryCode)
}

func (g *GeoIP) IsComplete() bool {
	for i := 0; i < (len(g.ipToCountryCode) - 1); i++ {
		currRange := g.ipToCountryCode[i]
		nextRange := g.ipToCountryCode[i+1]
		if (currRange.end + 1) != nextRange.start {
			return false
		}
	}
	return true
}

func (g *GeoIP) FindCountryByIP(ip string) (cc map[string]struct{}) {
	IPv4Address := net.ParseIP(ip).To4()
	if IPv4Address == nil {
		return
	}
	cc = make(map[string]struct{})
	mask := net.CIDRMask(*ipMask, defaultSubnetMask)
	isSubnet := IPv4Address.Equal(IPv4Address.Mask(mask))

	offset := 0
	if *ipMask != defaultSubnetMask && isSubnet {
		offset = wildcard
	}

	ipInt := IP4toInt(IPv4Address)
	for _, ipRange := range g.ipToCountryCode {
		if ipInt>>int64(offset) >= ipRange.start>>int64(offset) && ipInt>>int64(offset) <= ipRange.end>>int64(offset) {
			cc[ipRange.cc] = struct{}{}
		} else if len(cc) > 0 {
			return
		}
	}

	return
}

func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address)
	return IPv4Int.Int64()
}

func NewGeoIP() *GeoIP {
	// TODO: Do not forget to sort ip ranges ascending
	geoIP := &GeoIP{}

	// open file
	f, err := os.Open(*geoDB)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		startIP := rec[*ipStartIndex]
		endIP := rec[*ipEndIndex]
		countryCode := rec[*countryIndex]

		IPv4AddressStart := net.ParseIP(startIP).To4()
		IPv4AddressEnd := net.ParseIP(endIP).To4()
		if IPv4AddressStart == nil || IPv4AddressEnd == nil {
			continue
		}

		ipRange := IPRange{IP4toInt(IPv4AddressStart), IP4toInt(IPv4AddressEnd), countryCode}

		geoIP.ipToCountryCode = append(geoIP.ipToCountryCode, ipRange)
	}
	return geoIP
}

func findCountryCodes(geodb *GeoIP, f *os.File) (codes [][]string, err error) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		countryCodes := geodb.FindCountryByIP(input.Text())
		var cc []string
		for code := range countryCodes {
			cc = append(cc, code)
		}
		codes = append(codes, cc)
	}
	if input.Err() != nil {
		return nil, input.Err()
	}
	return
}

func handleFlags(geoip *GeoIP) {
	var quit bool
	if *showDbSize {
		fmt.Println("Total number of records are", geoip.Size())
		quit = true
	}
	if *checkDB {
		fmt.Println("GeoIP database covers all IPv4 addresses:", geoip.IsComplete())
		quit = true
	}
	if !(*ipMask > 0 && *ipMask <= 32) {
		fmt.Println("Subnet mask shall be whitin 1 - 32 range")
		quit = true
	}
	wildcard = defaultSubnetMask - *ipMask
	if quit {
		os.Exit(0)
	}
}

func main() {
	geoip := NewGeoIP()

	flag.Parse()
	handleFlags(geoip)

	var (
		codes [][]string
		err   error
	)

	files := flag.Args()
	if len(files) == 0 {
		codes, err = findCountryCodes(geoip, os.Stdin)
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, code := range codes {
		printCodes(code)
	}
}

func printCodes(codes []string) {
	for i, v := range codes {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
