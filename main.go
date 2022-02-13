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
	"strings"
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

func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address)
	return IPv4Int.Int64()
}

type GeoIP struct {
	GeoMap   map[int64]string
	GeoSlice []int64
}

func NewGeoIP(filename string, ipIndex, countryIndex int) *GeoIP {
	geoIP := &GeoIP{}
	geoIP.GeoMap = make(map[int64]string)

	// open file
	f, err := os.Open(filename)
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

		IPv4Address := net.ParseIP(rec[ipIndex]).To4()
		if IPv4Address == nil {
			continue
		}

		ipInt := IP4toInt(IPv4Address)
		geoIP.GeoMap[ipInt] = rec[countryIndex]
		geoIP.GeoSlice = append(geoIP.GeoSlice, ipInt)
	}
	return geoIP
}

func (g *GeoIP) FindCountry(ip net.IP) string {
	decIP := IP4toInt(ip)
	var prevD int64
	for _, d := range g.GeoSlice {
		if decIP == d {
			return g.GeoMap[d]
		} else if decIP < d {
			return g.GeoMap[prevD]
		}
		prevD = d
	}
	return ""
}

func cliLoop(geoIP *GeoIP) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter IPv4 address: ")
		ipn, _ := reader.ReadString('\n')

		ip := strings.Split(ipn, "\n")[0]

		IPv4Address := net.ParseIP(ip).To4()
		if IPv4Address == nil {
			fmt.Println("It is not an IPv4 address. Try again.")
			continue
		}

		fmt.Printf("IP address country: %q\n", geoIP.FindCountry(IPv4Address))
	}
}

func main() {

	flag.Parse()

	NewGeoIP(*geoDB, *ipIndex, *countryIndex)

}
