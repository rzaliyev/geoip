package main

import (
	"reflect"
	"testing"
)

const testDBFile = "test_data.csv"
const testDBFile2 = "test_data2.csv"

func TestGeoIP(t *testing.T) {

	t.Run("successfull reading of geodb", func(t *testing.T) {
		*geoDB = testDBFile
		geoip := NewGeoIP()

		want := 20
		got := geoip.Size()

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("successful find country by ip", func(t *testing.T) {
		*geoDB = testDBFile
		geoip := NewGeoIP()

		cases := []struct {
			ip      string
			country map[string]struct{}
		}{
			{"0.0.10.10", map[string]struct{}{"ZZ": {}}},
			{"1.0.0.200", map[string]struct{}{"AU": {}}},
			{"1.0.1.0", map[string]struct{}{"CN": {}}},
			{"1.0.25.20", map[string]struct{}{"JP": {}}},
			{"1.0.200.1", map[string]struct{}{"TH": {}}},
			{"10.15.200.17", map[string]struct{}{"ZZ": {}}},
			{"87.242.127.255", map[string]struct{}{"RU": {}}},
			{"127.0.0.1", map[string]struct{}{"ZZ": {}}},
			{"223.255.255.35", map[string]struct{}{"AU": {}}},
			{"255.255.255.255", map[string]struct{}{"ZZ": {}}},
		}

		for _, test := range cases {
			want := test.country
			got := geoip.FindCountryByIP(test.ip)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %q, want %q", got, want)
			}
		}
	})

	t.Run("successful find countries by ip subnet", func(t *testing.T) {
		*geoDB = testDBFile2
		*ipMask = 24
		wildcard = defaultSubnetMask - *ipMask
		geoip := NewGeoIP()

		cases := []struct {
			ip      string
			country map[string]struct{}
		}{
			{"0.0.0.0", map[string]struct{}{"ZZ": {}, "RU": {}, "UA": {}, "KZ": {}}},
			{"0.0.2.0", map[string]struct{}{"GB": {}, "US": {}}},
			{"0.0.2.200", map[string]struct{}{"US": {}}},
			{"0.0.3.0", map[string]struct{}{"FR": {}}},
		}

		for _, test := range cases {
			want := test.country
			got := geoip.FindCountryByIP(test.ip)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %q, want %q", got, want)
			}
		}
	})

	t.Run("verify inclomplete database", func(t *testing.T) {
		*geoDB = testDBFile
		geoip := NewGeoIP()

		want := false
		got := geoip.IsComplete()

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})

	t.Run("verify complete database", func(t *testing.T) {
		*geoDB = testDBFile2
		geoip := NewGeoIP()

		want := true
		got := geoip.IsComplete()

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})
}
