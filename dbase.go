package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math"
	"os"
	"strconv"
	"time"
)

// URL for connection to database
var configDB = os.Getenv("URL")

// Sensor struct describes basic object for work with database
type Sensor struct {
	ID          int       `gorm:"column:id; serial"`
	IDSensor    int       `gorm:"column:id_sensor; type:integer; not null"`
	ValueSensor float64   `gorm:"column:value_sensor; type float(2); not null"`
	TimeAdd     time.Time `gorm:"column:time_add; type:timestamp; not null"`
}

// TableName returns new name for table
func (Sensor) TableName() string {
	return "fict_sensors_syn"
}

// Avg struct describes object for receiving average value for different time ranges
type Avg struct {
	Avg float64 `gorm:"column:avg"`
}

// variables for values from database
var lastValue Sensor
var avgArr []float64
var avgValue Avg

// dbPostData insert new data into table
func dbPostData(idSens int, valueSens float64, ctx *gin.Context) {
	db, err := gorm.Open("postgres", configDB)
	if err != nil {
		ErrorResp(ctx, err.Error())
		return
	}

	defer func() {
		flag := db.Close()
		if flag != nil && err == nil {
			ErrorResp(ctx, err.Error())
			return
		}
	}()

	sensor := Sensor{IDSensor: idSens, ValueSensor: valueSens, TimeAdd: time.Now()}
	db.Create(&sensor)

	ctx.JSON(200, gin.H{
		"ErrorMSG": ""})
}

// dbGet provides connection to database and connects each kind of request with appropriate query
func dbGet(date string, ctx *gin.Context) {

	db, err := gorm.Open("postgres", configDB)
	if err != nil {
		ErrorResp(ctx, err.Error())
		return
	}

	defer func() {
		flag := db.Close()
		if flag != nil && err == nil {
			ErrorResp(ctx, err.Error())
		}
	}()

	idSensRaw := ctx.Query("id")
	idSens, _ := strconv.Atoi(idSensRaw)

	switch date {
	case "last":
		dbGetLast(idSens, db, ctx)
	case "day":
		dbGetDay(idSens, db, ctx)
	case "week":
		dbGetWeek(idSens, db, ctx)
	case "month":
		dbGetMonth(idSens, db, ctx)
	case "year":
		dbGetYear(idSens, db, ctx)

	// this default value is unreachable, but it's worth to leave it in case future architecture change
	default:
		ctx.JSON(404, gin.H{
			"ErrorMSG": "404"})
	}
}

// dbGetLast realizes query for the last value of certain sensor
func dbGetLast(idSens int, db *gorm.DB, ctx *gin.Context) {

	defer resetObjects()

	db.Where(Sensor{IDSensor: idSens}).Order("id desc").Limit(1).First(&lastValue)

	ctx.JSON(200, gin.H{
		"ErrorMSG": "", "values": math.Round(lastValue.ValueSensor*100) / 100})
}

// dbGetDay realizes query for average value for each of the last 24 hours
func dbGetDay(idSens int, db *gorm.DB, ctx *gin.Context) {

	defer resetObjects()

	// Okay, I know, this query is pretty far from being the best one, but there's no opportunity to use group by,
	// cause group by doesn't count hours/days/months where were no values and there is no more comfortable way
	// to count average value for each of hours. Although I was about to change it to complex query with subqueries
	// inside like this one and realize it let's say via loop:
	// SELECT (SELECT AVG(value_sensor) AS "avg"
	//         FROM sensors
	//         where id_sensor = 1
	//           and time_add >= now() - '2 day'::INTERVAL
	//           and time_add <= now() - '1 day'::INTERVAL) as value1,
	//        (SELECT AVG(value_sensor) AS "avg"
	//         FROM sensors
	//         where id_sensor = 1
	//           and time_add >= now() - '3 day'::INTERVAL
	//           and time_add <= now() - '2 day'::INTERVAL) as value2, ...;
	//
	// ... but then I realized that things are going messy and no one will understand this code after some time,
	// so as far as there are no thousands of users and database is on the same server as backend (at least this query)
	// things are not that bad as they seems to be at first
	// ::author @dedifferentiator
	for i := 0; i < 24; i++ {
		db.Raw(`SELECT AVG(value_sensor) AS "avg" FROM fict_sensors_syn where id_sensor=? and time_add >= now() 
					- ?::INTERVAL and time_add <= now() - ?::INTERVAL`, idSens, strconv.Itoa(i+1)+" hour",
			strconv.Itoa(i)+" hour").Scan(&avgValue)

		avgArr = append(avgArr, math.Round(avgValue.Avg*100)/100)
		avgValue.Avg = 0
	}

	ctx.JSON(200, gin.H{
		"ErrorMSG": "", "values": avgArr})
}

// dbGetWeek realizes query for average value for each of the last 7 days
func dbGetWeek(idSens int, db *gorm.DB, ctx *gin.Context) {

	defer resetObjects()

	for i := 0; i < 7; i++ {
		db.Raw(`SELECT AVG(value_sensor) AS "avg" FROM fict_sensors_syn where id_sensor=? and time_add >= now() 
					- ?::INTERVAL and time_add <= now() - ?::INTERVAL`, idSens, strconv.Itoa(i+1)+" day",
			strconv.Itoa(i)+" day").Scan(&avgValue)

		avgArr = append(avgArr, math.Round(avgValue.Avg*100)/100)
		avgValue.Avg = 0
	}

	ctx.JSON(200, gin.H{
		"ErrorMSG": "", "values": avgArr})
}

// dbGetMonth realizes query for average value for each of the last 30 days
func dbGetMonth(idSens int, db *gorm.DB, ctx *gin.Context) {

	defer resetObjects()

	for i := 0; i < 30; i++ {
		db.Raw(`SELECT AVG(value_sensor) AS "avg" FROM fict_sensors_syn where id_sensor=? and time_add >= now() - ?::INTERVAL
					and time_add <= now() - ?::INTERVAL`, idSens, strconv.Itoa(i+1)+" day",
			strconv.Itoa(i)+" day").Scan(&avgValue)

		avgArr = append(avgArr, math.Round(avgValue.Avg*100)/100)
		avgValue.Avg = 0
	}

	ctx.JSON(200, gin.H{
		"ErrorMSG": "", "values": avgArr})
}

// dbGetYear realizes query for average value for each of the last 12 months
func dbGetYear(idSens int, db *gorm.DB, ctx *gin.Context) {

	defer resetObjects()

	for i := 0; i < 12; i++ {
		db.Raw(`SELECT AVG(value_sensor) AS "avg" FROM fict_sensors_syn where id_sensor=? and time_add >= now() - ?::INTERVAL
					and time_add <= now() - ?::INTERVAL`, idSens, strconv.Itoa(i+1)+" month",
			strconv.Itoa(i)+" month").Scan(&avgValue)

		avgArr = append(avgArr, math.Round(avgValue.Avg*100)/100)
		avgValue.Avg = 0
	}

	ctx.JSON(200, gin.H{
		"ErrorMSG": "", "values": avgArr})
}

// resetObjects clear variables in dbGet... functions since they (variables) are in the outer scope and may not be
// overwritten due to a zero query result
func resetObjects() {
	lastValue.ValueSensor = 0
	avgArr = []float64{}
}
