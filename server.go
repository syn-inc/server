package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"strconv"
	"strings"
)

// portName defines port which app should open
var portName = ":" + os.Getenv("PORT")

// set defines path for postData request
var set = os.Getenv("SET")

// configRouter explores request for its HTTP-method and redirect it to appropriate function
func configRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("last", getPeriod)
	r.GET("day", getPeriod)
	r.GET("week", getPeriod)
	r.GET("month", getPeriod)
	r.GET("year", getPeriod)
	r.POST(set, postData)
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"ErrorMSG": "404"})
	})

	return r
}

// main runs server on selected port
func main() {
	r := configRouter()
	err := r.Run(portName)

	if err != nil {
		panic("Cannot open port " + portName)
	}
}

// postData tests validness of request
func postData(ctx *gin.Context) {

	idSensRaw := ctx.Query("id")
	valueSensRaw := ctx.Query("value")

	if IsSetOk(idSensRaw, valueSensRaw, ctx) {

		// err is muted due to validness check IsSetOk, so if there's an error it'll will be handled during check
		idSens, _ := strconv.Atoi(idSensRaw)
		valueSens, _ := strconv.ParseFloat(valueSensRaw, 64)
		dbPostData(idSens, valueSens, ctx)
	} else {
		ErrorResp(ctx, "Incorrect params")
	}
}

// getPeriod  validness of request
func getPeriod(ctx *gin.Context) {
	if IsGetOk(ctx) {
		// [1:] omits backslash
		dbGet(ctx.Request.URL.Path[1:], ctx)
	} else {
		ErrorResp(ctx, "Incorrect params")
	}
}

// IsSetOk Checks set-request for its correctness
func IsSetOk(idSens, valueSens string, ctx *gin.Context) (false bool) {

	// values which idSens and valueSens should not contain
	falseStruct := []string{"Inf", "NaN"}

	for _, value := range falseStruct {
		if strings.Contains(idSens, value) || strings.Contains(valueSens, value) {
			return
		}
	}

	keyVal, err := strconv.Atoi(idSens)
	if err != nil || keyVal <= 0 {
		return
	}

	_, err = strconv.ParseFloat(valueSens, 64)

	if err != nil {
		return
	}
	return true
}

// IsGetOk test get-request
func IsGetOk(ctx *gin.Context) bool {

	idSens := ctx.Query("id")

	if strings.Contains(idSens, "Inf") || strings.Contains(idSens, "NaN") {
		return false
	}

	keyVal, err := strconv.Atoi(idSens)

	if err != nil || keyVal <= 0 {
		return false
	}
	return true
}

// ErrorResp return JSON in response body with 500 code as result of wrong request and error message which describes it
func ErrorResp(ctx *gin.Context, err string) {
	ctx.JSON(500, gin.H{"ErrorMSG": err})
}
