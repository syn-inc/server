package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
