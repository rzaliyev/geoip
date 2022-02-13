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
			country string
		}{
			{"0.0.10.10", "ZZ"},
			{"1.0.0.200", "AU"},
			{"1.0.1.0", "CN"},
			{"1.0.25.20", "JP"},
			{"1.0.200.1", "TH"},
			{"10.15.200.17", "ZZ"},
			{"87.242.127.255", "RU"},
			{"127.0.0.1", "ZZ"},
			{"223.255.255.35", "AU"},
			{"255.255.255.255", "ZZ"},
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

	t.Run("verify clomplete database", func(t *testing.T) {
		*geoDB = testDBFile2
		geoip := NewGeoIP()

		want := true
		got := geoip.IsComplete()

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})
}
