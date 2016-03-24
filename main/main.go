package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Toubib/elmo"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/mreiferson/go-httpclient"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	VERSION    = "0.1-dev"
	BUILD_DATE = ""
)

//cli flags
var (
	mainUrl        = flag.String("url", "", "The url to get.")
	version        = flag.Bool("version", false, "Print version information.")
	verbose        = flag.Bool("verbose", false, "Print more informations.")
	parallelFetch  = flag.Int("parallel", 8, "Number of parallel fetch to launch. 0 means unlimited.")
	connectTimeout = flag.Int("connect-timeout", 1000, "Connect timeout in ms.")

	useNagios          = flag.Bool("use-nagios", false, "Nagios compatible output.")
	nagiosWarningTime  = flag.Int("nagios-warning", 5000, "Nagios warning time in ms.")
	nagiosCriticalTime = flag.Int("nagios-critical", 10000, "Nagios critical time in ms.")

	requestTimeout        = flag.Int("request-timeout", 10000, "Request timeout in ms.")
	responseHeaderTimeout = flag.Int("response-header-timeout", 0, "Response header timeout in ms.")
	useInflux             = flag.Bool("use-influx", false, "Send data to influxdb.")
	influxUrl             = flag.String("influx-url", "http://localhost:8086", "The influx database access url.")
	influxDatabase        = flag.String("influx-database", "elmo", "The influx database name.")
	assetsAllowedDomains  = flag.String("assets-allowed-domains", "", "List of allowed assets domains to fetch from, comma separated.")
)

func main() {

	flag.Parse()

	if *version {
		fmt.Printf("%v\nBuild: %v\n", VERSION, BUILD_DATE)
		return
	}

	o := elmo.Options(*verbose)

	elmo.run(o)

	// send data to influxdb
	if *useInflux {
		sendstatsToInflux(&assetsStats)
	}

	// We're done! Print the results...
	if *useNagios {
		fmt.Printf("Downloaded %vKB in %d/%d files in %v.|size=%vKB time=%v;%v;%v;0;%v\n",
			gstat.totalResponseSize/1024, len(assetsStats), len(assets), gstat.totalResponseTime,
			gstat.totalResponseSize/1024, gstat.totalResponseTime, *nagiosWarningTime, *nagiosCriticalTime, *requestTimeout,
		)

		//nagios exit
		if gstat.totalResponseTime >= time.Duration(*nagiosCriticalTime)*time.Millisecond {
			os.Exit(NAGIOS_ERROR)
		} else if gstat.totalResponseTime >= time.Duration(*nagiosWarningTime)*time.Millisecond {
			os.Exit(NAGIOS_WARNING)
		} else {
			os.Exit(NAGIOS_OK)
		}

	} else {
		fmt.Printf("Downloaded assets: %d/%d.\n", len(assetsStats), len(assets))
		fmt.Printf("Total time: %v.\n", cyan(gstat.totalResponseTime))
		fmt.Printf("Total size: %v%s.\n", white(gstat.totalResponseSize/1024), white("kb"))
	}
}
