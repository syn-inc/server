package main

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// URL for connection to test database, should be the same as configDB
var configTestDB = os.Getenv("URLTEST")

// TestMain setup environment for testing
func TestMain(m *testing.M) {

	db, err := gorm.Open("postgres", configTestDB)
	if err != nil {
		panic("Cannot open connection to test database")
	}

	defer func() {
		flag := db.Close()
		if flag != nil && err == nil {
			panic("Cannot close connection to test database")
		}
	}()

	db.Exec(`CREATE TABLE fict_sensors_syn(id serial primary key, id_sensor integer not null, value_sensor float(2)
				 not null, time_add timestamp not null); INSERT INTO fict_sensors_syn(id_sensor, value_sensor,
 				 time_add) VALUES (1, 1, now()), (1, 2, now() - '1 day'::INTERVAL),
												 (1, 3, now() - '2 day'::INTERVAL),
												 (1, 4, now() - '3 day'::INTERVAL),
												 (1, 5, now() - '4 day'::INTERVAL),
												 (2, 10, now() - '2 day'::INTERVAL),
												 (2, 11, now() - '3 day'::INTERVAL),
												 (2, 12, now() - '4 day'::INTERVAL),
												 (2, 13, now() - '5 day'::INTERVAL),
												 (2, 14, now() - '6 day'::INTERVAL),
												 (2, 15, now() - '7 day'::INTERVAL);`)
	testCode := m.Run()
	db.Exec(`DROP TABLE fict_sensors_syn;`)
	os.Exit(testCode)
}

// TestGetWeek1 test GET-week request
func TestGetWeek1(t *testing.T) {

	router := configRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/week?id=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[1,2,3,4,5,0,0]}`, w.Body.String())
}

// TestGetWeek2 test GET-week request
func TestGetWeek2(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/week?id=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[0,0,10,11,12,13,14]}`, w.Body.String())
}

// TestGetWeek3 test GET-week request
func TestGetWeek3(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/week?id=3", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[0,0,0,0,0,0,0]}`, w.Body.String())
}
