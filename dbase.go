package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

const (
	host     = os.Getenv("HOST")
	port     = 5432
	user     = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbName   = os.Getenv("DATABASE")
)

func Main2(sensId int) {
	GetLastValue(sensId)
	GetLastDay(sensId)
	fmt.Println("--------------------------------")
	GetLastWeek(sensId)
}

func GetLastValue(sensId int) {
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
	var idValue string
	rows, err := db.Query(`SELECT value_sensor FROM "538" where id_sensor=$1 order by id DESC LIMIT 1`, sensId)
	if err != nil {
		panic(err)
	}

	defer func() {
		flag := rows.Close()
		if flag != nil && err == nil {
			panic(err)
		}
	}()
	for rows.Next() {
		err = rows.Scan(&idValue)
		if err != nil {
			panic(err)
		}
		fmt.Println(idValue)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
}

func GetLastDay(sensId int) {

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
	var idValue string
	var hourValue string
	rows, err := db.Query(`SELECT EXTRACT(hour from time_add), value_sensor FROM "538" where time_add >= NOW() - '1 day'::INTERVAL and id_sensor=$1`, sensId)
	if err != nil {
		panic(err)
	}

	defer func() {
		flag := rows.Close()
		if flag != nil && err == nil {
			panic(err)
		}
	}()
	for rows.Next() {
		err = rows.Scan(&hourValue, &idValue)
		if err != nil {
			panic(err)
		}
		fmt.Println(hourValue, idValue)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
}
func GetLastWeek(sensId int) {

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

	var idValue string
	var hourValue string
	rows, err := db.Query(`SELECT EXTRACT(hour from time_add), value_sensor FROM "538" where time_add >= NOW() - '1 week'::INTERVAL and id_sensor=$1`, sensId)
	if err != nil {
		panic(err)
	}

	defer func() {
		flag := rows.Close()
		if flag != nil && err == nil {
			panic(err)
		}
	}()

	for rows.Next() {
		err = rows.Scan(&hourValue, &idValue)
		if err != nil {
			panic(err)
		}
		fmt.Println(hourValue, idValue)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
}
