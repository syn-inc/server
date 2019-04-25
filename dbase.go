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
var configStr = os.Getenv("URL")

// Sensor struct describes basic object for work with database
type Sensor struct {
	Id          int       `gorm:"column:id; serial"`
	IdSensor    int       `gorm:"column:id_sensor; type:integer; not null"`
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

var lastValue Sensor
var avgArr []float64
var avgValue Avg

// dbPostData insert new data into table
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

	sensor := Sensor{IdSensor: idSens, ValueSensor: valueSens, TimeAdd: time.Now()}
	db.Create(&sensor)

	ctx.JSON(200, gin.H{
		"ErrorMSG": ""})
}

// dbGet provides connection to database and connects each kind of request with appropriate query
func dbGet(date string, ctx *gin.Context) {

	db, err := gorm.Open("postgres", configStr)
	if err != nil {
		ErrorResp(ctx, err.Error())
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
	default:
		ctx.JSON(404, gin.H{
			"ErrorMSG": ""})
	}
}

// dbGetLast realizes query for the last value of certain sensor
func dbGetLast(idSens int, db *gorm.DB, ctx *gin.Context) {

	defer resetObjects()

	db.Raw(`SELECT value_sensor FROM fict_sensors_syn where id_sensor=? order by id desc limit 1;`,
		idSens).Scan(&lastValue)
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

		// Firstly it looks like there's a bug here, avgValue.Avg is not clearing each iteration, but it's not a bug,
		// it is a feature! The resulting array will be used later for charts, so instead of breaking smooth line of
		// curve it just repeats the last non-null value if it exist. Agree, maybe it's not the best solution, but for
		// now we decided to leave things as they are, although they'll might be changed later.
		//
		avgArr = append(avgArr, math.Round(avgValue.Avg*100)/100)
	}

	ctx.JSON(200, gin.H{
		"ErrorMSG": "", "values": avgArr})
}

// dbGetWeek realizes query for average value for each of the last 7 days
func dbGetWeek(idSens int, db *gorm.DB, ctx *gin.Context) {

	defer resetObjects()

	for i := 0; i < 7; i++ {
		db.Raw(`SELECT AVG(value_sensor) AS "avg" FROM fict_sensors_syn where id_sensor=? and time_add >= now() - ?::INTERVAL
					and time_add <= now() - ?::INTERVAL`, idSens, strconv.Itoa(i+1)+" day",
			strconv.Itoa(i)+" day").Scan(&avgValue)

		avgArr = append(avgArr, math.Round(avgValue.Avg*100)/100)
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
	}

	ctx.JSON(200, gin.H{
		"ErrorMSG": "", "values": avgArr})
}

// dbGetYear realizes query for average value for each of the last 12 month
func dbGetYear(idSens int, db *gorm.DB, ctx *gin.Context) {

	defer resetObjects()

	for i := 0; i < 12; i++ {
		db.Raw(`SELECT AVG(value_sensor) AS "avg" FROM fict_sensors_syn where id_sensor=? and time_add >= now() - ?::INTERVAL
					and time_add <= now() - ?::INTERVAL`, idSens, strconv.Itoa(i+1)+" month",
			strconv.Itoa(i)+" month").Scan(&avgValue)

		avgArr = append(avgArr, math.Round(avgValue.Avg*100)/100)
	}

	ctx.JSON(200, gin.H{
		"ErrorMSG": "", "values": avgArr})
}

// resetObjects clear variables in dbGet... functions since they (variables) are in the outer scope and may not be
// overwritten due to the lack of a non-zero query result
func resetObjects() {
	lastValue.ValueSensor = 0
	avgValue.Avg = 0
	avgArr = []float64{}
}
