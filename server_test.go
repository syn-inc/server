package main

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
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

// TestGetLast tests GET-last request
func TestGetLast(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/last?id=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":3.6}`, w.Body.String())
}

// TestGetWeek tests GET-day request
func TestGetDay(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/day?id=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[1,0,2.5,0,4.5,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,4,0]}`,
		w.Body.String())
}

// TestGetWeek tests GET-week request
func TestGetWeek(t *testing.T) {

	router := configRouter()

	comparison := map[string]string{"/week?id=1": `{"ErrorMSG":"","values":[3.17,0,0,4.5,5,0,0]}`,
		"/week?id=2": `{"ErrorMSG":"","values":[0,0,10.13,11,50.02,0,14.01]}`,
		"/week?id=3": `{"ErrorMSG":"","values":[0,0,0,0,0,0,0]}`}

	for key, value := range comparison {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", key, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, value, w.Body.String())
	}
}

// TestGetMonth tests GET-month request
func TestGetMonth(t *testing.T) {
	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/month?id=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"ErrorMSG":"","values":[0,0,10.13,11,50.02,0,14.01,15.67,0,0,0,0,0,0,11.5,0,0,0,0,0,`+
		`13.5,0,0,0,0,0,0,0,0,0]}`, w.Body.String())
}

// TestGetYear tests GET-year request
func TestGetYear(t *testing.T) {
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
	assert.Equal(t, 0, sensValue.ID)
	assert.Equal(t, 1, sensValue.IDSensor)
	assert.Equal(t, 21.01, math.Round(sensValue.ValueSensor*100)/100)
}

// TestError404 tests wrong path for GET/POST-request
func TestError404(t *testing.T) {
	router := configRouter()

	methods := []string{"GET", "POST"}
	requests := []string{"/fff?id=2", "/Last?id=2", "/YEAR?id=2", "/111?id=2", "/?id=2", "/lastweek?id=2"}

	for _, method := range methods {
		for _, request := range requests {

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, request, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, 404, w.Code)
			assert.Equal(t, `{"ErrorMSG":"404"}`, w.Body.String())
		}
	}
}

// TestError500 tests lack of arguments for GET/POST-request
func TestError500(t *testing.T) {
	router := configRouter()

	comparison := map[string]string{"GET": "/last", "POST": "/" + set}

	for key, value := range comparison {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(key, value, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, `{"ErrorMSG":"Incorrect params"}`, w.Body.String())
	}
}

// TestDbPostData tests dbPostData func
func TestDbPostData(t *testing.T) {
	configDB = "host=localhost port=5433 user=postgres dbname=WRONG-NAME password=PASSWORD sslmode=disable"

	router := configRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/"+set+"?id=1&value=20.06", nil)
	router.ServeHTTP(w, req)

	type Resp struct {
		ErrorMSG string
	}

	var responseRes Resp
	// decoding json from response body to check that there is no error in ErrorMSG
	err := json.NewDecoder(w.Body).Decode(&responseRes)
	if err != nil {
		panic("Cannot decode JSON from body of response")
	}

	assert.Equal(t, 500, w.Code)
	assert.NotEqual(t, "", responseRes.ErrorMSG)
}

// TestDbGetData tests dbPostData func with wrong connection URL
func TestDbGet(t *testing.T) {
	configDB = "host=localhost port=5433 user=postgres dbname=WRONG-NAME password=PASSWORD sslmode=disable"

	router := configRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/last?id=1", nil)
	router.ServeHTTP(w, req)

	type Resp struct {
		ErrorMSG string
	}
	var responseRes Resp
	// decoding json from response body to check that there is no error in ErrorMSG
	err := json.NewDecoder(w.Body).Decode(&responseRes)
	if err != nil {
		panic("Cannot decode JSON from body of response")
	}
	assert.Equal(t, 500, w.Code)
	assert.NotEqual(t, "", responseRes.ErrorMSG)
}

// TestIsGetOk tests IsGetOk func with wrong arguments in request
func TestIsGetOk(t *testing.T) {
	router := configRouter()

	requests := [...]string{"/last?id=k", "/last?id=Inf", "/last?id=NaN", "/last?id=-1", "/last?id=0", "/last?id=1.2"}

	for _, val := range requests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", val, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, `{"ErrorMSG":"Incorrect params"}`, w.Body.String())
	}
}

// TestIsSetOk tests IsSetOk func with wrong arguments in request
func TestIsSetOk(t *testing.T) {
	router := configRouter()

	requests := [...]string{"/" + set + "?id=1&value=Inf", "/" + set + "?id=-1&value=23", "/" + set + "?id=&value=1",
		"/" + set + "?id=Inf&value=1", "/" + set + "?id=1&value=", "/" + set + "?id=Inf&value=", "/" + set +
		"?id=NaN&value=1", "/" + set + "?id=1&value=NaN", "/" + set + "?id=NaN&value=NaN", "/" + set +
		"?id=NaN&value=NaN"}

	for _, val := range requests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", val, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, `{"ErrorMSG":"Incorrect params"}`, w.Body.String())
	}
}
