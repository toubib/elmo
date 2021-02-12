# elmo ![Go](https://github.com/toubib/elmo/workflows/Go/badge.svg)

Elmo let you fetch a web page with its assets.

It can send download time statistics to [InfluxDB](https://www.influxdata.com/products/influxdb/) and be used as a [Nagios](https://www.nagios.org) check plugin.

! This project is not released yet, use with caution !

## Usage

```
Usage of ./elmo:
  -assets-allowed-domains string
    	List of allowed assets domains to fetch from, comma separated.
  -connect-timeout int
    	Connect timeout in ms. (default 1000)
  -influx-database string
    	The influx database name. (default "elmo")
  -influx-url string
    	The influx database access url. (default "http://localhost:8086")
  -nagios-critical int
    	Nagios critical time in ms. (default 10000)
  -nagios-warning int
    	Nagios warning time in ms. (default 5000)
  -parallel int
    	Number of parallel fetch to launch. 0 means unlimited. (default 8)
  -request-timeout int
    	Request timeout in ms. (default 10000)
  -response-header-timeout int
    	Response header timeout in ms.
  -url string
    	The url to get.
  -use-influx
    	Send data to influxdb.
  -use-nagios
    	Nagios compatible output.
  -verbose
    	Print more informations.
  -version
    	Print version information.
```
