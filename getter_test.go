package main

import (
	"net/url"
	"testing"
)

type CheckTest []bool

func TestGetter(t *testing.T) {
	IsSetOk(url.Values{"temp": {"56"}, "hum": {"24"}})
	IsSetOk(url.Values{"temp": {"-56"}, "hum": {"-24"}, "hi": {"-0"}})
	IsSetOk(url.Values{"temp": {"0"}, "hum": {"0"}})
	notTruthStructure := CheckTest{
		IsSetOk(url.Values{}),
		IsSetOk(url.Values{"temp": {}, "hum": {}}),
		IsSetOk(url.Values{"temp": {""}, "hum": {""}}),
		IsSetOk(url.Values{"": {"54"}, "_": {"67"}}),
		IsSetOk(url.Values{"Temp": {"56"}, "hUm": {"24"}}),
		IsSetOk(url.Values{"temp": {"ab"}, "hum": {"24"}}),
		IsSetOk(url.Values{"temp": {"56"}, "hum": {"cd"}}),
		IsSetOk(url.Values{"temp": {"_"}, "hum": {"_"}}),
		IsSetOk(url.Values{"temp": {"nil"}, "hum": {"nil"}}),
		IsSetOk(url.Values{"temp": {"inf"}, "hum": {"inf"}}),
		IsSetOk(url.Values{"temp": {"e"}, "hum": {"e"}}),
		IsSetOk(url.Values{"temp": {"56"}, "hum": {"24"}, "hi": {"0"}}),
		IsSetOk(url.Values{"temp": {"56"}, "hum": {"24"}, "hi": {"_0"}}),
		IsSetOk(url.Values{"temp": {"56"}, "hum": {"2-4"}, "hi": {"0"}})}
	for _, value := range notTruthStructure {
		if value == true {
			t.Log("Failed test")
			break
		}
	}
}
