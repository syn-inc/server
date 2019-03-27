package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	host     = os.Getenv("HOST")
	port     = 5432
	user     = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbName   = os.Getenv("DATABASE")
)

func ServerBody(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatalf("Parse Error %s", err)
	}

	fmt.Println(r.Form)
	fmt.Println("Path", r.URL.Path)
	fmt.Println("Scheme", r.URL.Scheme)
	ViewShow(w, "Hello World!")
	if r.URL.Path == "/get" {
		getRequest(w)
	} else if r.URL.Path == "/set" { //&& IsSetOk(r.Form) == true
		TempHumSet(w, r.Form)
	}
}

func main() {
	http.HandleFunc("/", ServerBody)
	err := http.ListenAndServe("", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func dbSet(idSens int, sensValue float64) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer func() {
		flag := db.Close()
		if flag != nil && err == nil {
			panic(err)
		}
	}()

	// should be used for checking and establishment connection to db, but which way?
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var idValue int
	sqlGetLastId := `SELECT MAX(id) FROM "538"`
	responseLastId := db.QueryRow(sqlGetLastId)
	err = responseLastId.Scan(&idValue)
	if err != nil {
		panic(err)
	}

	sqlStatement := `INSERT INTO "538" (id, id_sensor, value_sensor, time_add) VALUES ($1, $2, $3, to_timestamp($4, 'yyyy-mm-dd hh24:mi:ss'))`
	db.QueryRow(sqlStatement, idValue+1, idSens, sensValue, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Setting db error")
	}
}

func TempHumSet(w http.ResponseWriter, Form url.Values) {

	for key, value := range Form {
		flagValue, err := strconv.ParseFloat(value[0], 64)
		if err != nil {
			panic(err)
		}

		fmt.Println(key, flagValue)
		if flagValue != 0 {
			newKey, err := strconv.Atoi(key)
			if err != nil {
				panic(err)
			}
			dbSet(newKey, flagValue)
			ViewShow(w, "\nSuccessfully set!")
		}
		flagValue = 0
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
func getRequest(w http.ResponseWriter) {

	PSQLInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := sql.Open("postgres", PSQLInfo)
	if err != nil {
		panic(err)
	}

	defer func() {
		flag := db.Close()
		if flag != nil && err == nil {
			panic(err)
		}
	}()

	var idValue int
	err = db.QueryRow(`SELECT MAX(id) FROM "538"`).Scan(&idValue)
	if err != nil {
		panic("Max value error")
	}

	var idSensor int
	var valueSensor float64
	var timeAdd string
	err = db.QueryRow(`SELECT id_sensor, value_sensor, time_add FROM "538" where id=$1`, idValue).Scan(&idSensor, &valueSensor, &timeAdd)
	if err != nil {
		panic(err)
	}
	fmt.Println(timeAdd)
	ViewShow(w, "\nid: "+strconv.Itoa(idSensor)+"\nvalue: "+strconv.FormatFloat(valueSensor, 'f', 2, 64)+"\ndate: "+timeAdd)
}
