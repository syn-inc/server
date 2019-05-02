package main

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"math"
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
 				 time_add) VALUES (1, 1, now()),
                                  (1, 2, now() - '2 hour'::INTERVAL),
							      (1, 3, now() - '2 hour'::INTERVAL),
								  (1, 4, now() - '4 hour'::INTERVAL),
								  (1, 5, now() - '4 hour'::INTERVAL),
								  (1, 4, now() - '22 hour'::INTERVAL),
								  (1, 5, now() - '3 day'::INTERVAL),
								  (1, 4, now() - '3 day'::INTERVAL),
								  (1, 5, now() - '4 day'::INTERVAL),
 
								  (2, 10, now() - '2 day'::INTERVAL),
								  (2, 10.25, now() - '2 day'::INTERVAL),
								  (2, 11, now() - '3 day'::INTERVAL),
								  (2, 13.99, now() - '4 day'::INTERVAL),
								  (2, 95.56, now() - '4 day'::INTERVAL),
								  (2, 60.35, now() - '4 day'::INTERVAL),
								  (2, 30.20, now() - '4 day'::INTERVAL),
								  (2, 14.01, now() - '6 day'::INTERVAL),

								  (2, 15.67, now() - '7 day'::INTERVAL),
								  (2, 11.00, now() - '14 day'::INTERVAL),
								  (2, 12, now() - '14 day'::INTERVAL),
								  (2, 13, now() - '20 day'::INTERVAL),
								  (2, 14, now() - '20 day'::INTERVAL),

								  (2, 50, now() - '1 month'::INTERVAL),
								  (2, 79, now() - '2 month'::INTERVAL),
								  (2, 45.96, now() - '2 month'::INTERVAL),
								  (2, 1, now() - '4 month'::INTERVAL),
								  (2, 99.9999999999, now() - '4 month'::INTERVAL),
								  (2, 87.1, now() - '6 month'::INTERVAL),
								  (2, 52.94, now() - '10 month'::INTERVAL),
								  (2, 5.85, now() - '10 month'::INTERVAL),
								  (2, 3.6, now() - '11 month'::INTERVAL);`)
	testCode := m.Run()
	//drop table after testing
	db.Exec(`DROP TABLE fict_sensors_syn;`)
	os.Exit(testCode)
}

// TestGetLast1 tests GET-last request
func TestGetLast1(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/last?id=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":3.6}`, w.Body.String())
}

// TestGetWeek1 tests GET-day request
func TestGetDay1(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/day?id=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[1,0,2.5,0,4.5,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,4,0]}`, w.Body.String())
}

// TestGetWeek1 tests GET-week request
func TestGetWeek1(t *testing.T) {

	router := configRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/week?id=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[3.17,0,0,4.5,5,0,0]}`, w.Body.String())
}

// TestGetWeek2 tests GET-week request
func TestGetWeek2(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/week?id=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[0,0,10.13,11,50.02,0,14.01]}`, w.Body.String())
}

// TestGetWeek3 tests GET-week request
func TestGetWeek3(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/week?id=3", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[0,0,0,0,0,0,0]}`, w.Body.String())
}

// TestGetMonth1 tests GET-month request
func TestGetMonth1(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/month?id=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[0,0,10.13,11,50.02,0,14.01,15.67,0,0,0,0,0,0,11.5,0,0,0,0,0,`+
		`13.5,0,0,0,0,0,0,0,0,0]}`, w.Body.String())
}

// TestGetYear1 test GET-year request
func TestGetYear1(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/year?id=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[23.93,50,62.48,0,50.5,0,87.1,0,0,0,29.39,3.6]}`, w.Body.String())
}

// TestPost tests dbPostData request
func TestPost(t *testing.T) {

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

	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/"+set+"?id=1&value=21.01", nil)
	router.ServeHTTP(w, req)

	var sensValue Sensor
	db.Raw(`SELECT id_sensor, value_sensor FROM fict_sensors_syn WHERE id_sensor=1 ORDER BY id DESC LIMIT 
			     1`).Scan(&sensValue)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":""}`, w.Body.String())
	assert.Equal(t, 0, sensValue.Id)
	assert.Equal(t, 1, sensValue.IdSensor)
	assert.Equal(t, 21.01, math.Round(sensValue.ValueSensor*100)/100)
}

// TestError404N1 tests wrong path for GET-request
func TestError404N1(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/fff?id=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, `{"ErrorMSG":"404"}`, w.Body.String())
}

// TestError404N2 tests wrong path for POST-request
func TestError404N2(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/fff?id=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, `{"ErrorMSG":"404"}`, w.Body.String())
}

// TestError404N1 tests lack of arguments for GET-request
func TestError500N1(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/last", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"ErrorMSG":"Incorrect params"}`, w.Body.String())
}

// TestError500N2 tests lack of arguments for POST-request
func TestError500N2(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/"+set, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"ErrorMSG":"Incorrect params"}`, w.Body.String())
}

// TestDbPostData testing dbPostData func
func TestDbPostData(t *testing.T) {
	configDB = "host=localhost port=5433 user=postgres dbname=WRONG-NAME password=PASSWORD sslmode=disable"

	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/"+set+"?id=1&value=20.06", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.NotEqual(t, `{"ErrorMSG":""}`, w.Body.String())
}
