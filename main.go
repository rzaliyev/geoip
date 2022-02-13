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

const (
	defaultInputFilename            = "data.csv"
	defaultInputIpAddressStartIndex = 0
	defaultInputIpAddressEndIndex   = 1
	defaultInputCountryIndex        = 2
)

var (
	geoDB        = flag.String("geodb", defaultInputFilename, "csv file with geoip data")
	ipStartIndex = flag.Int("ipstart", defaultInputIpAddressStartIndex, "ip address start index in csv file (default 0)")
	ipEndIndex   = flag.Int("ipend", defaultInputIpAddressEndIndex, "ip address end index in csv file")
	countryIndex = flag.Int("country", defaultInputCountryIndex, "country index in csv file")
)

type IPIntRange struct {
	start int64
	end   int64
}

type IPStringRange struct {
	start string
	end   string
}

type GeoIP struct {
	ipToCountryCode       map[IPIntRange]string
	countryCodeToIPRanges map[string][]IPStringRange
}

func (g *GeoIP) Size() int {
	return len(g.ipToCountryCode)
}

func (g *GeoIP) FindCountryByIP(ip string) string {
	IPv4Address := net.ParseIP(ip).To4()
	if IPv4Address == nil {
		return ""
	}

	ipInt := IP4toInt(IPv4Address)
	for ipRange, cc := range g.ipToCountryCode {
		if ipInt >= ipRange.start && ipInt <= ipRange.end {
			return cc
		}
	}
	return ""
}

func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address)
	return IPv4Int.Int64()
}

func NewGeoIP() *GeoIP {
	geoIP := &GeoIP{}
	geoIP.ipToCountryCode = make(map[IPIntRange]string)
	geoIP.countryCodeToIPRanges = make(map[string][]IPStringRange)

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

		IPv4AddressStart := net.ParseIP(startIP).To4()
		IPv4AddressEnd := net.ParseIP(endIP).To4()
		if IPv4AddressStart == nil || IPv4AddressEnd == nil {
			continue
		}
		countryCode := rec[*countryIndex]
		strRange := IPStringRange{startIP, endIP}
		intRange := IPIntRange{IP4toInt(IPv4AddressStart), IP4toInt(IPv4AddressEnd)}

		geoIP.ipToCountryCode[intRange] = countryCode
		geoIP.countryCodeToIPRanges[countryCode] = append(geoIP.countryCodeToIPRanges[countryCode], strRange)
	}
	return geoIP
}

func findCountryCodes(geodb *GeoIP, f *os.File) (codes []string, err error) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		codes = append(codes, geodb.FindCountryByIP(input.Text()))
	}
	if input.Err() != nil {
		return nil, input.Err()
	}
	return
}

func main() {

	flag.Parse()

	geoip := NewGeoIP()
	fmt.Println(geoip.Size())

	var (
		codes []string
		err   error
	)

	files := os.Args[1:]
	if len(files) == 0 {
		codes, err = findCountryCodes(geoip, os.Stdin)
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, code := range codes {
		fmt.Println(code)
	}

}
