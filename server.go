package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"strconv"
	"strings"
)

var portName = ":" + os.Getenv("PORT")
var set = os.Getenv("SET")

// configRouter explores request for its HTTP-method and redirect it to appropriate function
func configRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/last", getLast)
	r.GET("/day", getDay)
	r.GET("/week", getWeek)
	r.GET("/month", getMonth)
	r.GET("/year", getYear)
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

// postData test validness of request
func postData(ctx *gin.Context) {

	idSensRaw := ctx.Query("id")
	valueSensRaw := ctx.Query("value")

	if IsSetOk(idSensRaw, valueSensRaw, ctx) {

		// err is muted due to validness check IsSetOk, so if there's an error it'll will be handled during check
		idSens, _ := strconv.Atoi(idSensRaw)
		valueSens, _ := strconv.ParseFloat(valueSensRaw, 64)
		dbPostData(idSens, valueSens, ctx)
	}
}

// getLast test validness of request
func getLast(ctx *gin.Context) {
	if IsGetOk(ctx) {
		dbGet("last", ctx)
	}
}

// getDay test validness of request
func getDay(ctx *gin.Context) {
	if IsGetOk(ctx) {
		dbGet("day", ctx)
	}
}

// getWeek test validness of request
func getWeek(ctx *gin.Context) {
	if IsGetOk(ctx) {
		dbGet("week", ctx)
	}
}

// getMonth test validness of request
func getMonth(ctx *gin.Context) {
	if IsGetOk(ctx) {
		dbGet("month", ctx)
	}
}

// getYear test validness of request
func getYear(ctx *gin.Context) {
	if IsGetOk(ctx) {
		dbGet("year", ctx)
	}
}

// IsSetOk Checks set-request for its correctness
func IsSetOk(idSens, valueSens string, ctx *gin.Context) (false bool) {

	falseStruct := []string{"Inf", "NaN"}

	for _, value := range falseStruct {
		if strings.Contains(idSens, value) || strings.Contains(valueSens, value) {
			ErrorResp(ctx, "Incorrect params")
			return
		}
	}

	keyVal, err := strconv.Atoi(idSens)
	if err != nil || keyVal <= 0 {
		ErrorResp(ctx, "Incorrect params")
		return
	}

	_, err = strconv.ParseFloat(valueSens, 64)

	if err != nil {
		ErrorResp(ctx, "Incorrect params")
		return
	}
	return true
}

// IsGetOk test get-request
func IsGetOk(ctx *gin.Context) bool {

	idSens := ctx.Query("id")

	if strings.Contains(idSens, "Inf") || strings.Contains(idSens, "NaN") {
		ErrorResp(ctx, "Incorrect params")
		return false
	}

	keyVal, err := strconv.Atoi(idSens)

	if err != nil || keyVal <= 0 {
		ErrorResp(ctx, "Incorrect params")
		return false
	}
	return true
}

// ErrorResp return JSON in response body with 500 code as result of wrong request and error message which describes it
func ErrorResp(ctx *gin.Context, err string) {
	ctx.JSON(500, gin.H{"ErrorMSG": err})
}
