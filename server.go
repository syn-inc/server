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

func ServerBody(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatalf("Parse Error %s", err)
	}

	if r.URL.Path == "/set" && IsSetOk(r.Form) == true {
		SetData(w, r.Form)
	} else if r.URL.Path == "/get" && IsGetOk(r.Form) {
		GetData(w, r.Form)
	}
	//TODO add return json with request error
}

func main() {
	http.HandleFunc("/", ServerBody)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func IsSetOk(v url.Values) bool {
	if len(v) != 1 {
		return false
	}

	for key, value := range v {

		if len(value) != 1 {
			return false
		}

		if strings.Contains(value[0], "Inf") || strings.Contains(value[0], "NaN") {
			return false
		}

		keyVal, err := strconv.Atoi(key)
		if err != nil {
			return false
		}

		if keyVal <= 0 {
			return false
		}

		_, err = strconv.ParseFloat(value[0], 64)
		if err != nil {
			return false
		}
	}
	return true
}

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
			if optionsList[value[0]] == false {
				return false
			}
		}
	}
	return true
}

func SetData(w http.ResponseWriter, form url.Values) {

	for key, value := range form {
		newValue, err := strconv.ParseFloat(value[0], 64)
		if err != nil {
			panic(err)
		}
		newKey, err := strconv.Atoi(key)
		if err != nil {
			panic(err)
		}
		dbSet(newKey, newValue)
		ViewShow(w, "\nSuccessfully set!")
	}
}

func ViewShow(w http.ResponseWriter, s string) {
	_, err := fmt.Fprintf(w, s)
	if err != nil {
		log.Fatalf("Parse Error %s", err)
	}
}

func GetData(w http.ResponseWriter, form url.Values) {

	w.Header().Set("Content-Type", "application/json")

	var err error
	var idSens int
	var date string

	for key, value := range form {
		// TODO write tests
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

	type LastValues struct {
		ErrorMsg string
		Values   []float64
	}

	lastValues := LastValues{"IncorrectParams", []float64{}}
	switch date {
	case "last":
		lastValues = LastValues{"ok", dbGetLastValue(idSens)}
	case "day":
		lastValues = LastValues{"ok", dbGetLastDay(idSens)}
	case "week":
		lastValues = LastValues{"ok", dbGetLastWeek(idSens)}
	case "month":
		lastValues = LastValues{"ok", dbGetLastMonth(idSens)}
	case "year":
		lastValues = LastValues{"ok", dbGetLastYear(idSens)}
	}
	jsonObj, err := json.Marshal(lastValues)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = w.Write(jsonObj)
	if err != nil {
		panic("WriteError")
	}
}
