package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//Structure for default JSON-file
type JsonResp struct {
	ErrorMsg string
	Values   []float64
}

//ServerBody process requests
func ServerBody(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatalf("Parse Error %s", err)
	}

	if r.URL.Path == os.Getenv("SET") && IsSetOk(r.Form) {
		SetData(w, r.Form)
	} else if r.URL.Path == "/get" && IsGetOk(r.Form) {
		GetData(w, r.Form)
	} else {

		w.Header().Set("Content-Type", "application/json")
		response := JsonResp{"404", []float64{}}
		jsonObj, err := json.Marshal(response)

		if err != nil {
			panic("Cannot Marshal 404 JSON")
		}

		_, err = w.Write(jsonObj)
		if err != nil {
			panic("Cannot write 404 JSON")
		}
	}
}

//main run server for defined port
func main() {
	http.HandleFunc("/", ServerBody)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//IsSetOk Checks set-request for its correctness
func IsSetOk(v url.Values) bool {

	if len(v) != 2 {
		return false
	}

	if len(v["id"]) != 1 || len(v["value"]) != 1 {
		return false
	}

	for key, value := range v {

		if len(value) != 1 {
			return false
		}

		if strings.Contains(value[0], "Inf") || strings.Contains(value[0], "NaN") {
			return false
		}

		if key == "id" {
			keyVal, err := strconv.Atoi(value[0])
			if err != nil {
				return false
			}
			if keyVal <= 0 {
				return false
			}
		}

		if key == "value" {
			_, err := strconv.ParseFloat(value[0], 64)
			if err != nil {
				return false
			}
		}
	}
	return true
}

//IsGetOk Checks get-request for its correctness
func IsGetOk(v url.Values) bool {

	if len(v) != 2 {
		return false
	}

	if len(v["id"]) != 1 || len(v["date"]) != 1 {
		return false
	}

	for key, value := range v {

		if len(value) != 1 {
			return false
		}

		if strings.Contains(value[0], "Inf") || strings.Contains(value[0], "NaN") {
			return false
		}

		if key == "id" {
			keyVal, err := strconv.Atoi(value[0])
			if err != nil {
				return false
			}
			if keyVal <= 0 {
				return false
			}
		}

		if key == "date" {
			optionsList := map[string]bool{"last": true, "day": true, "week": true, "month": true, "year": true}
			if !optionsList[value[0]] {
				return false
			}
		}
	}
	return true
}

//SetData Processes given data and transmit it to func dbSet()
func SetData(w http.ResponseWriter, form url.Values) {

	var newKey int
	var newValue float64

	for key, value := range form {
		switch key {
		case "id":
			newKey, _ = strconv.Atoi(value[0])
		case "value":
			newValue, _ = strconv.ParseFloat(value[0], 64)
		}
	}
	dbSet(newKey, newValue)
	ViewShow(w, "Successfully set!")
}

//ViewShowOutput given string on the page
func ViewShow(w http.ResponseWriter, s string) {
	_, err := fmt.Fprintf(w, s)
	if err != nil {
		log.Fatalf("Parse Error %s", err)
	}
}

//GetData Processes and output given data then transmit it to func dbGetLast...()
func GetData(w http.ResponseWriter, form url.Values) {

	w.Header().Set("Content-Type", "application/json")

	var err error
	var idSens int
	var date string

	//Parse values and keys from get request
	for key, value := range form {
		switch key {
		case "id":
			idSens, err = strconv.Atoi(value[0])
			if err != nil {
				ViewShow(w, "ID parse error, please contact developer to fix such errors")
			}
		case "date":
			date = value[0]
		}
	}

	//Creating appropriate to request JSON-file
	lastValues := JsonResp{"IncorrectParams", []float64{}}
	switch date {
	case "last":
		lastValues = JsonResp{"ok", dbGetLastValue(idSens)}
	case "day":
		lastValues = JsonResp{"ok", dbGetLastDay(idSens)}
	case "week":
		lastValues = JsonResp{"ok", dbGetLastWeek(idSens)}
	case "month":
		lastValues = JsonResp{"ok", dbGetLastMonth(idSens)}
	case "year":
		lastValues = JsonResp{"ok", dbGetLastYear(idSens)}
	}
	jsonObj, err := json.Marshal(lastValues)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//Output JSON
	_, err = w.Write(jsonObj)
	if err != nil {
		panic("WriteError")
	}
}
