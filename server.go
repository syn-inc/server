package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strconv"
	"strings"
)

func main() {
	r := gin.Default()
	r.GET("/get", getLast)
	r.POST("", postData)

	portName := ":8000"
	err := r.Run(portName)
	if err != nil {
		panic("Cannot open port " + portName)
	}
}

func getLast(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong"})
}

func postData(ctx *gin.Context) {
	idSensRaw := ctx.Query("id")
	valueSensRaw := ctx.Query("value")
	idSens, _ := strconv.Atoi(idSensRaw)
	valueSens, _ := strconv.ParseFloat(valueSensRaw, 64)
	if IsSetOk(idSensRaw, valueSensRaw) {
		dbPostData(idSens, valueSens, ctx)
	}
}

//IsSetOk Checks set-request for its correctness
func IsSetOk(idSens, valueSens string) bool {

	if strings.Contains(idSens, "Inf") || strings.Contains(idSens, "NaN") {
		return false
	}

	if strings.Contains(valueSens, "Inf") || strings.Contains(valueSens, "NaN") {
		return false
	}

	keyVal, err := strconv.Atoi(idSens)
	if err != nil {
		return false
	}
	if keyVal <= 0 {
		return false
	}

	_, err = strconv.ParseFloat(valueSens, 64)
	if err != nil {
		return false
	}
	return true
}
