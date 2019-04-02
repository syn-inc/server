package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func ServerBody(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatalf("Parse Error %s", err)
	}

	if r.URL.Path == "/set" { //&& IsSetOk(r.Form) == true
		SetData(w, r.Form)
	} else if r.URL.Path == "/get" {
		getRequest(w, r.Form)
	}
}

func main() {
	http.HandleFunc("/", ServerBody)
	err := http.ListenAndServe(":"+"8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func SetData(w http.ResponseWriter, form url.Values) {

	for key, value := range form {
		newValue, err := strconv.ParseFloat(value[0], 64)
		if err != nil {
			panic(err)
		}

		fmt.Println(key, newValue)
		if newValue != 0 {
			newKey, err := strconv.Atoi(key)
			if err != nil {
				panic(err)
			}
			dbSet("538", newKey, newValue)
			// FIXME add processing of error
			ViewShow(w, "\nSuccessfully set!")
		}
		newValue = 0
	}
}

func ViewShow(w http.ResponseWriter, s string) {
	_, err := fmt.Fprintf(w, s)
	if err != nil {
		log.Fatalf("Parse Error %s", err)
	}
}

func IsSetOk(v url.Values) bool { //t *testing.T,
	if len(v) != 2 {
		//t.Log("Number of keys is not equal 2, ", v)
		return false
	}
	if len(v["temp"]) != 1 && len(v["hum"]) != 1 {
		//t.Log("Number of arguments should be equal 1")
		return false
	}
	_, err := strconv.ParseFloat(v["temp"][0], 64)
	if err != nil {
		//t.Log("Temp's argument isn't Float64 ", v)
		return false
	}
	_, err = strconv.ParseFloat(v["hum"][0], 64)
	if err != nil {
		//t.Log("Humidity's argument isn't Float64, ", v)
		return false
	}
	return true
}

func getRequest(w http.ResponseWriter, form url.Values) {

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
				jsonObj, err := json.Marshal(map[string]string{"errorMsg": "IncorrectParams"})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, err = w.Write(jsonObj)
				if err != nil {
					panic(err)
				}
			}
		case "date":
			date = value[0]
		}
	}

	type LastValues struct {
		ErrorMsg string
		Values   []float64
	}

	switch date {
	case "last":
		lastValues := LastValues{"ok", dbGetLastValue(idSens)}
		jsonObj, err := json.Marshal(lastValues)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonObj)
		if err != nil {
			panic("Wrong Write")
		}
	case "day":
		lastValues := LastValues{"ok", dbGetLastDay(idSens)}
		jsonObj, err := json.Marshal(lastValues)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonObj)
		if err != nil {
			panic(err)
		}
	case "week":
		lastValues := LastValues{"ok", dbGetLastWeek(idSens)}
		jsonObj, err := json.Marshal(lastValues)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonObj)
		if err != nil {
			panic(err)
		}
	case "month":

		lastValues := LastValues{"ok", dbGetLastMonth(idSens)}
		jsonObj, err := json.Marshal(lastValues)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonObj)
		if err != nil {
			panic(err)
		}
	case "year":
		lastValues := LastValues{"ok", dbGetLastYear(idSens)}
		jsonObj, err := json.Marshal(lastValues)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonObj)
		if err != nil {
			panic(err)
		}
	default:
		jsonObj, err := json.Marshal(map[string]string{"errorMsg": "IncorrectParams"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonObj)
		if err != nil {
			panic(err)
		}
	}
}
