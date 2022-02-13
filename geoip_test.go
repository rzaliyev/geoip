package main

import (
	"testing"
)

const testDBFile = "test_data.csv"

func TestGeoIP(t *testing.T) {

	t.Run("successfull reading of geodb", func(t *testing.T) {
		*geoDB = testDBFile
		geoip := NewGeoIP()

		want := 10
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
		}

		for _, test := range cases {
			want := test.country
			got := geoip.FindCountry(test.ip)
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		}
	})
}
