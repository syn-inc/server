package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

var configStr = os.Getenv("URL")

type Sensor struct {
	Id           int       `gorm:"serial"`
	Id_sensor    int       `gorm:"type:integer; not null"`
	Value_sensor float64   `gorm:"type float(2); not null"`
	Time_add     time.Time `gorm:"type:timestamp; not null"`
}

func dbPostData(idSens int, valueSens float64, ctx *gin.Context) {
	db, err := gorm.Open("postgres", configStr)
	if err != nil {
		ctx.JSON(500, gin.H{"ErrorMSG": err.Error()})
		panic(err)
	}

	defer func() {
		flag := db.Close()
		if flag != nil && err == nil {
			panic(err)
		}
	}()

	sensor := Sensor{Id_sensor: idSens, Value_sensor: valueSens, Time_add: time.Now()}
	db.Create(&sensor)

	ctx.JSON(200, gin.H{
		"ErrorMSG": ""})
}
