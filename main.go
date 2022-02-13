package main

import (
	"encoding/csv"
	"flag"
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

func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address)
	return IPv4Int.Int64()
}

type GeoIP struct {
	countries  []string
	ipIntegers []int64
}

func (g *GeoIP) Size() int {
	return len(g.ipIntegers)
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

// func cliLoop(geoIP *GeoIP) {
// 	for {
// 		reader := bufio.NewReader(os.Stdin)
// 		fmt.Print("Enter IPv4 address: ")
// 		ipn, _ := reader.ReadString('\n')

// 		ip := strings.Split(ipn, "\n")[0]

// 		IPv4Address := net.ParseIP(ip).To4()
// 		if IPv4Address == nil {
// 			fmt.Println("It is not an IPv4 address. Try again.")
// 			continue
// 		}

// 		fmt.Printf("IP address country: %q\n", geoIP.FindCountry(IPv4Address))
// 	}
// }

func main() {

	flag.Parse()

	NewGeoIP()

}
