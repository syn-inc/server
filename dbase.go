package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math"
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

// this method should only be used on tables with id column
func dbSet(tableName string, idSens int, sensValue float64) {
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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var idValue int
	//noinspection SqlResolve
	err = db.QueryRow(`SELECT MAX(id) FROM $1`, tableName).Scan(&idValue)
	if err != nil {
		panic(err)
	}

	sqlStatement := `INSERT INTO $1 VALUES ($2, $3, $4, to_timestamp($5, 'yyyy-mm-dd hh24:mi:ss'))`
	db.QueryRow(sqlStatement, tableName, idValue+1, idSens, sensValue, time.Now().Format("2000-01-01 00:00:00"))
	if err != nil {
		log.Println("Setting db error")
	}
}

func dbGet() (int, float64, string) {

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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var idValue int
	err = db.QueryRow(`SELECT MAX(id) FROM "538"`).Scan(&idValue)
	if err != nil {
		panic(err)
	}

	var idSensor int
	var valueSensor float64
	var timeAdd string
	err = db.QueryRow(`SELECT id_sensor, value_sensor, time_add FROM "538" where id=$1`, idValue).Scan(&idSensor, &valueSensor, &timeAdd)
	if err != nil {
		panic(err)
	}
	return idSensor, valueSensor, timeAdd
}

func dbGetLastValue(idSens int) []float64 {
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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var value float64
	err = db.QueryRow(`SELECT value_sensor FROM "538" where id_sensor=$1 order by id DESC LIMIT 1`, idSens).Scan(&value)
	if err != nil {
		panic("Querying error")
	}

	// round value
	return []float64{math.Round(value*100) / 100}
}

func dbGetLastDay(sensId int) []float64 {

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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var value float64
	var valueArr []float64
	// FIXME nil error on parse value
	for i := 0; i < 24; i++ {
		err = db.QueryRow(`SELECT AVG(value_sensor) AS "Average value" FROM "538" where id_sensor=$1 and
                                time_add >= now() - $2::INTERVAL and time_add <= now() - $3::INTERVAL`, sensId, strconv.Itoa(i+1)+" hour", strconv.Itoa(i)+" hour").Scan(&value)
		if err != nil {
			value = 10000
		}
		valueArr = append(valueArr, math.Round(value*100)/100)
		err = nil
	}
	fmt.Println(valueArr)
	return valueArr
}

func dbGetLastWeek(sensId int) []float64 {

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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var value float64
	var valueArr []float64
	// FIXME nil error on parse value
	for i := 0; i < 7; i++ {
		err = db.QueryRow(`SELECT AVG(value_sensor) AS "Average value" FROM "538" where id_sensor=$1 and
                                time_add >= now() - $2::INTERVAL and time_add <= now() - $3::INTERVAL`, sensId,
			strconv.Itoa(i+1)+" day", strconv.Itoa(i)+" day").Scan(&value)
		if err != nil {
			value = 10000
		}
		valueArr = append(valueArr, math.Round(value*100)/100)
	}
	fmt.Println(valueArr)
	return valueArr
}

func dbGetLastMonth(sensId int) []float64 {

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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var value float64
	var valueArr []float64
	// FIXME nil error on parse value
	for i := 0; i < 30; i++ {
		err = db.QueryRow(`SELECT AVG(value_sensor) AS "Average value" FROM "538" where id_sensor=$1 and
                                time_add >= now() - $2::INTERVAL and time_add <= now() - $3::INTERVAL`, sensId, strconv.Itoa(i+1)+" day", strconv.Itoa(i)+" day").Scan(&value)
		if err != nil {
			value = 10000
		}
		valueArr = append(valueArr, math.Round(value*100)/100)
	}
	fmt.Println(valueArr)
	return valueArr
}

func dbGetLastYear(sensId int) []float64 {

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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var value float64
	var valueArr []float64
	// FIXME nil error on parse value
	for i := 0; i < 12; i++ {
		err = db.QueryRow(`SELECT AVG(value_sensor) AS "Average value" FROM "538" where id_sensor=$1 and
                                time_add >= now() - $2::INTERVAL and time_add <= now() - $3::INTERVAL and extract(year from now())=extract(year from time_add)`, sensId, strconv.Itoa(i+1)+" month", strconv.Itoa(i)+" month").Scan(&value)
		if err != nil {
			value = 10000
		}
		valueArr = append(valueArr, math.Round(value*100)/100)
	}
	fmt.Println(valueArr)
	return valueArr
}
