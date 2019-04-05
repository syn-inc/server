package main

import (
	"net/url"
	"testing"
)

type CheckTest []bool

func TestIsSetOk(t *testing.T) {
	truthStructure := CheckTest{
		IsGetOk(url.Values{"id": {"1"}, "value": {"2.05"}}),
		IsGetOk(url.Values{"id": {"541"}, "value": {"-1"}}),
		IsGetOk(url.Values{"id": {"4"}, "value": {"0"}}),
		IsGetOk(url.Values{"id": {"421455"}, "value": {"0.9999999999999"}}),
		IsGetOk(url.Values{"id": {"67890523"}, "value": {"0.01"}}),
		IsGetOk(url.Values{"id": {"100"}, "value": {"6"}}),
	}

	notTruthStructure := CheckTest{
		IsSetOk(url.Values{}),
		IsSetOk(url.Values{"": {}}),
		IsSetOk(url.Values{"1": {""}}),
		IsSetOk(url.Values{"1": {""}, "hum": {""}}),
		IsSetOk(url.Values{"fdf": {"1"}, "hum": {""}}),
		IsSetOk(url.Values{"fdf": {"2,e2e"}}),
		IsSetOk(url.Values{"fdf": {"2"}}),
		IsSetOk(url.Values{"1": {"g"}}),
		IsSetOk(url.Values{"1": {"nil"}}),
		IsSetOk(url.Values{"1": {"nil"}}),
		IsSetOk(url.Values{"-1": {"nil"}}),
		IsSetOk(url.Values{"-1": {"23.02"}}),
		IsSetOk(url.Values{"0.5": {"23.02"}}),
		IsSetOk(url.Values{"0.9999999999999999999999999": {"23.02"}}),
		IsSetOk(url.Values{"21453.9999999999999999999999999": {"0.02"}}),
		IsSetOk(url.Values{"-1": {"nil"}}),
		IsSetOk(url.Values{"101": {"Inf"}}),
		IsSetOk(url.Values{"302": {"-Inf"}}),
		IsSetOk(url.Values{"324536245643224214": {"NaN"}}),

		IsGetOk(url.Values{"date": {"1"}, "id": {"last"}}),
		IsGetOk(url.Values{"vf": {"v"}, "v": {"vff"}}),
		IsGetOk(url.Values{"id": {""}, "date": {""}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {""}}),
		IsGetOk(url.Values{"id": {""}, "date": {"month"}}),
		IsGetOk(url.Values{"id": {"-1"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"0"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"Last"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"LAST"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"lastlast"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"lastmonth"}}),
		IsGetOk(url.Values{"id": {"0.25"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"0.999999999999999999999999"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {""}}),
		IsGetOk(url.Values{"id": {"NaN"}, "date": {"year"}}),
		IsGetOk(url.Values{"id": {"Inf"}, "date": {"day"}}),
		IsGetOk(url.Values{"id": {"-Inf"}, "date": {"day"}}),
		IsGetOk(url.Values{"id": {""}, "date": {"day"}}),
	}
	for _, value := range truthStructure {
		if value == false {
			t.Fatal("Failed IsSetOk, truth part")
		}
	}
	for _, value := range notTruthStructure {
		if value == true {
			t.Fatal("Failed IsSetOk, notTruth part")
		}
	}
}

func TestIsGetOk(t *testing.T) {
	truthStructure := CheckTest{
		IsGetOk(url.Values{"id": {"1"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"day"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"week"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"month"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"year"}}),
		IsGetOk(url.Values{"id": {"14"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"53263246"}, "date": {"month"}}),
	}

	notTruthStructure := CheckTest{
		IsGetOk(url.Values{}),
		IsGetOk(url.Values{"": {}}),
		IsGetOk(url.Values{"1": {""}}),
		IsGetOk(url.Values{"1": {""}, "hum": {""}}),
		IsGetOk(url.Values{"fdf": {"1"}, "hum": {""}}),
		IsGetOk(url.Values{"fdf": {"2,e2e"}}),
		IsGetOk(url.Values{"fdf": {"2"}}),
		IsGetOk(url.Values{"1": {"g"}}),
		IsGetOk(url.Values{"1": {"nil"}}),
		IsGetOk(url.Values{"1": {"nil"}}),
		IsGetOk(url.Values{"-1": {"nil"}}),
		IsGetOk(url.Values{"-1": {"23.02"}}),
		IsGetOk(url.Values{"-1": {"nil"}}),
		IsGetOk(url.Values{"101": {"Inf"}}),
		IsGetOk(url.Values{"302": {"-Inf"}}),
		IsGetOk(url.Values{"324536245643224214": {"NaN"}}),

		IsGetOk(url.Values{"date": {"1"}, "id": {"last"}}),
		IsGetOk(url.Values{"vf": {"v"}, "v": {"vff"}}),
		IsGetOk(url.Values{"id": {""}, "date": {""}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {""}}),
		IsGetOk(url.Values{"id": {""}, "date": {"month"}}),
		IsGetOk(url.Values{"id": {"-1"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"0"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"Last"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"LAST"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"lastlast"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {"lastmonth"}}),
		IsGetOk(url.Values{"id": {"0.25"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"0.999999999999999999999999"}, "date": {"last"}}),
		IsGetOk(url.Values{"id": {"1"}, "date": {""}}),
		IsGetOk(url.Values{"id": {"NaN"}, "date": {"year"}}),
		IsGetOk(url.Values{"id": {"Inf"}, "date": {"day"}}),
		IsGetOk(url.Values{"id": {"-Inf"}, "date": {"day"}}),
		IsGetOk(url.Values{"id": {""}, "date": {"day"}}),
	}
	for _, value := range truthStructure {
		if value == false {
			t.Fatal("Failed IsGetOk, truth part")
		}
	}
	for _, value := range notTruthStructure {
		if value == true {
			t.Fatal("Failed IsGetOk, notTruth part")
		}
	}
}
