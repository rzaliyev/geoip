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
	defaultInputFilename       = "data.csv"
	defaultInputIpAddressIndex = 0
	defaultInputCountryIndex   = 2
)

var (
	geoDB        = flag.String("geodb", defaultInputFilename, "csv file with geoip data")
	ipIndex      = flag.Int("ip", defaultInputIpAddressIndex, "ip address index in csv file (default 0)")
	countryIndex = flag.Int("country", defaultInputCountryIndex, "country index in csv file")
)

type GeoIP struct {
	countries  []string
	ipIntegers []int64
}

func (g *GeoIP) Size() int {
	return len(g.ipIntegers)
}

func (g *GeoIP) FindCountry(ip string) string {
	IPv4Address := net.ParseIP(ip).To4()
	if IPv4Address == nil {
		return ""
	}

	decIP := IP4toInt(IPv4Address)
	var prev int
	for i, d := range g.ipIntegers {
		if decIP == d {
			return g.countries[i]
		} else if decIP < d {
			return g.countries[prev]
		}
		prev = i
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

		IPv4Address := net.ParseIP(rec[*ipIndex]).To4()
		if IPv4Address == nil {
			continue
		}

		geoIP.countries = append(geoIP.countries, rec[*countryIndex])
		geoIP.ipIntegers = append(geoIP.ipIntegers, IP4toInt(IPv4Address))
	}
	return geoIP
}

func findCountryCodes(geodb *GeoIP, f *os.File) (codes []string, err error) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		codes = append(codes, geodb.FindCountry(input.Text()))
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
