[![Build Status](https://travis-ci.com/syn-inc/server.svg?branch=master)](https://travis-ci.com/syn-inc/server)
[![golangci](https://golangci.com/badges/github.com/syn-inc/server.svg)](https://golangci.com/r/github.com/syn-inc/server)
[![Go Report Card](https://goreportcard.com/badge/github.com/syn-inc/server)](https://goreportcard.com/report/github.com/syn-inc/server)
[![Maintainability](https://api.codeclimate.com/v1/badges/1fd3631ebfa1173067c2/maintainability)](https://codeclimate.com/github/syn-inc/server/maintainability)

# SYN-server
## Samples of requests:
### To set data:
To set data you have to give one parameter, key is `sensor_id` and value is its value. So the sample of request looks like this `/set?1=23.5`
### To get data
To get data you have to give two parameters, `id=$sensor_id` and `date`, `date` can be only any of these ones - `"last", "day", "week", "month", "year"`. Sample of request `/get?id=214432&date=week`.
- **last** - last value
- **day** - average values for each of the last 24 hours
- **week** - average values for each of the last 7 days
- **month** - average values for each of the last 30 days
- **year** - average values for each of the last 12 months
