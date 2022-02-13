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

const defaultInputFilename = "data.csv"

const (
	defaultInputIpAddressStartIndex = iota
	defaultInputIpAddressEndIndex
	defaultInputCountryIndex
)

var (
	geoDB        = flag.String("geodb", defaultInputFilename, "csv file with geoip data")
	ipStartIndex = flag.Int("ipstart", defaultInputIpAddressStartIndex, "ip address start index in csv file (default 0)")
	ipEndIndex   = flag.Int("ipend", defaultInputIpAddressEndIndex, "ip address end index in csv file")
	countryIndex = flag.Int("country", defaultInputCountryIndex, "country index in csv file")
	showDbSize   = flag.Bool("size", false, "show number of records in database")
	checkDB      = flag.Bool("check", false, "check geoip for completeness")
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

func (g *GeoIP) FindCountryByIP(ip string) string {
	IPv4Address := net.ParseIP(ip).To4()
	if IPv4Address == nil {
		return ""
	}

	ipInt := IP4toInt(IPv4Address)
	for _, ipRange := range g.ipToCountryCode {
		if ipInt >= ipRange.start && ipInt <= ipRange.end {
			return ipRange.cc
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

func hanleFlags(geoip *GeoIP) {
	flag.Parse()

	var quit bool

	if *showDbSize {
		fmt.Println("Total number of records are", geoip.Size())
		quit = true
	}
	if *checkDB {
		fmt.Println("GeoIP database covers all IPv4 addresses:", geoip.IsComplete())
		quit = true
	}

	if quit {
		os.Exit(0)
	}
}

func main() {

	geoip := NewGeoIP()
	hanleFlags(geoip)

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
