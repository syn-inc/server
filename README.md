<img align="right" width="250px" src="media/logo.png">

[![Build Status](https://travis-ci.com/syn-inc/server.svg?branch=master)](https://travis-ci.com/syn-inc/server)
[![golangci](https://golangci.com/badges/github.com/syn-inc/server.svg)](https://golangci.com/r/github.com/syn-inc/server)
[![Go Report Card](https://goreportcard.com/badge/github.com/syn-inc/server)](https://goreportcard.com/report/github.com/syn-inc/server)
[![Maintainability](https://api.codeclimate.com/v1/badges/1fd3631ebfa1173067c2/maintainability)](https://codeclimate.com/github/syn-inc/server/maintainability)
[![codecov](https://codecov.io/gh/syn-inc/server/branch/master/graph/badge.svg)](https://codecov.io/gh/syn-inc/server)



# SYN-server
## Samples of requests:
### To get data
Get data request consists of two parameters, path and key, e.g. `/last?id=1`. Path defines period and key - sensor id. 
Period might be only among these values - `"last", "day", "week", "month", "year"`.
- **last** - last value
- **day** - average values for each of the last 24 hours
- **week** - average values for each of the last 7 days
- **month** - average values for each of the last 30 days
- **year** - average values for each of the last 12 months

##Available sensors
 **ID - DECRYPTION**


- **1** - _Temperature, °C_
- **2** - _Humidity, %_
- **3** - _Pressure, hPa_
- **5** - _Amount of luminous flux, lm_

