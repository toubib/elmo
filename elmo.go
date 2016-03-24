/*
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>. 1
*/

package elmo

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
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

type Options struct {
	verbose bool
}

type downloadStatistic struct {
	url          string
	responseTime time.Duration
	responseSize int
	statusCode   int
}

type globalStatistic struct {
	totalResponseTime time.Duration
	totalResponseSize int
}

//colors !
var (
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	white  = color.New(color.FgWhite).SprintFunc()
	red    = color.New(color.FgRed, color.Bold).SprintFunc()
	green  = color.New(color.FgGreen, color.Bold).SprintFunc()
)

var globalStartTime = time.Now()

// Helper function to pull the  attribute from a Token
func getLink(t *html.Token) (ok bool, link string) {

	// Check link types, we need only stylesheet
	if t.Data == "link" {
		for _, a := range t.Attr {
			if a.Key == "rel" && a.Val != "stylesheet" {
				ok = false
				return
			}
		}

	}

	// search for css style assets on div
	if t.Data == "div" {
		ok = false
		for _, a := range t.Attr {

			// search for style key
			if a.Key == "style" {

				// search url
				if strings.Contains(a.Val, "url(") {
					r, _ := regexp.Compile(`url *\( *['"](.*)['"] *\)`)
					searchResult := r.FindStringSubmatch(a.Val)

					//url found
					if searchResult != nil {
						link = searchResult[1]
						ok = true
					}
					return
				}
			}
		}
		return
	}

	// Iterate over all of the Token's attributes until we find an "src"
	for _, a := range t.Attr {
		if a.Key == "src" || a.Key == "href" {
			link = a.Val
			ok = true
		}
	}

	return
}

// Extract all http** links from a given webpage
func fetchMainUrl(mainUrl *string, client *http.Client, verbose *bool) ([]string, downloadStatistic, error) {

	//List of urls found
	var assets []string

	//set downloadStatistic
	stat := downloadStatistic{*mainUrl, 0, 0, 0}

	//timer before
	t0 := time.Now()

	//launch the query
	req, _ := http.NewRequest("GET", *mainUrl, nil)
	resp, err := client.Do(req)

	if err != nil {
		return assets, stat, err
	}

	//Set stats
	stat.responseTime = time.Now().Sub(t0)
	stat.statusCode = resp.StatusCode

	//get the body size
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return assets, stat, err
	}

	//Set response size stat
	stat.responseSize = len(body)

	//Print download
	if *verbose {
		fmt.Printf("%s\t%s %s %v %v%s\n", time.Now().Sub(globalStartTime), green(stat.statusCode), stat.url, cyan(stat.responseTime), white(stat.responseSize), white("b"))
	}

	//extract assets from html
	assets = extractAssets(&body, req)

	return assets, stat, nil
}

//Get a html body and extract all assets links
func extractAssets(body *[]byte, mainRequest *http.Request) []string {
	var assets []string

	//create the tokenizer
	z := html.NewTokenizer(bytes.NewReader(*body))

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			return assets
		case html.SelfClosingTagToken, html.StartTagToken:
			t := z.Token()

			// Check if the token is a target tag
			// And extract the link if there is one
			if t.Data != "img" &&
				t.Data != "script" &&
				t.Data != "embed" &&
				t.Data != "div" && //only url in styles
				t.Data != "link" && //only stylesheet
				t.Data != "input" {

				continue
			}

			linkFound, assetUrl := getLink(&t)

			// We do not found a link
			if !linkFound {
				continue
			}
			//fmt.Println("link found:", assetUrl)

			// Make sure the url start with http
			if strings.Index(assetUrl, "http") != 0 {
				//get asset url object
				u, err := url.Parse(assetUrl)
				if err != nil {
					fmt.Println(red("ERROR -- "), err)
					continue
				}
				assetUrl = mainRequest.URL.ResolveReference(u).String()
			}
			assets = append(assets, assetUrl)
		}
	}
}

//Fetch an asset and get downloadStatistic
func fetchAsset(assetUrl string, client *http.Client, chStat chan downloadStatistic, chFinished chan bool, verbose *bool) {

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	//set downloadStatistic
	stat := downloadStatistic{assetUrl, 0, 0, 0}

	//timer before
	t0 := time.Now()

	//launch the query
	req, _ := http.NewRequest("GET", assetUrl, nil)

	if checkIfDomainAllowed(&req.URL.Host) == false {
		return
	}

	resp, err := client.Do(req)

	//handle error
	if err != nil {
		if !*useNagios {
			fmt.Println(red("Error:"), stat.url, err)
		}
		return
	}

	//Set stat
	stat.responseTime = time.Now().Sub(t0)
	stat.statusCode = resp.StatusCode

	//get the body size
	b := resp.Body
	defer b.Close() // close Body when the function returns
	body, err := ioutil.ReadAll(resp.Body)

	//Set response size stat
	stat.responseSize = len(body)

	//Print download
	if *verbose {
		fmt.Printf("%s\t%s %s %v %v%s\n", time.Now().Sub(globalStartTime), green(stat.statusCode), stat.url, cyan(stat.responseTime), white(stat.responseSize), white("b"))
	}

	chStat <- stat
}

//Send statistic data to influxdb
func sendstatsToInflux(assetsStats *[]downloadStatistic) {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: *influxUrl})
	if err != nil {
		fmt.Println(red("Influxdb - error creating InfluxDB Client:\n"), err.Error())
	}

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  *influxDatabase,
		Precision: "s",
	})

	//we want all the points of the batch to the same time
	influxTime := time.Now()

	//prepare data for influx
	for _, stat := range *assetsStats {
		// Create a point and add to batch
		tags := map[string]string{"url": stat.url}
		//var tags map[string]string
		fields := map[string]interface{}{
			"responseTime": int64(stat.responseTime),
			"responseSize": stat.responseSize,
		}
		pt, err := client.NewPoint(*mainUrl, tags, fields, influxTime)
		if err != nil {
			fmt.Println(red("Influxdb - error:\n"), err.Error())
		}
		bp.AddPoint(pt)
		//if (*verbose) {
		//	fmt.Printf("influx: %v\n", pt)
		//}
	}

	//Send data to influx
	err = c.Write(bp)
	if err != nil {
		fmt.Println(red("Influxdb - error:\n"), err.Error())
	}
}

// test if the given domain is allowed to fetch
func checkIfDomainAllowed(host *string) bool {

	if *assetsAllowedDomains == "" {
		return true
	}

	allowedDomains := strings.Split(*assetsAllowedDomains, ",")

	for _, domain := range allowedDomains {
		if domain == *host {
			return true
		}
	}

	return false
}

func run(o Options) globalStatistic {
	//urls and global stats
	var (
		assets      []string
		assetsStats []downloadStatistic
		mainUrlStat downloadStatistic
		//gstat           globalStatistic
		currentUrlIndex int
		err             error
	)

	// Channels
	chUrls := make(chan downloadStatistic)
	chFinished := make(chan bool)

	//set global timeout if nagios timeout is set
	if *useNagios {
		*requestTimeout = *nagiosCriticalTime
	}

	//Set a transport with timeouts
	transport := &httpclient.Transport{
		ConnectTimeout:        time.Duration(*connectTimeout) * time.Millisecond,
		RequestTimeout:        time.Duration(*requestTimeout) * time.Millisecond,
		ResponseHeaderTimeout: time.Duration(*responseHeaderTimeout) * time.Millisecond,
	}
	defer transport.Close()

	//Set an http client with this transport
	client := &http.Client{Transport: transport}

	//Set timer for global time
	t0 := time.Now()

	//Fetch the main url and get inner links
	assets, mainUrlStat, err = fetchMainUrl(mainUrl, client, &o.verbose)

	//handle main url error
	if err != nil {
		if *useNagios {
			fmt.Println(err)
			os.Exit(NAGIOS_ERROR)
		} else {
			fmt.Println(red("Fatal:"), err)
			os.Exit(1)
		}
	}

	//add main url response time
	gstat.totalResponseSize += mainUrlStat.responseSize

	//Fetch the firsts inner links
	for _, assetUrl := range assets {
		//fmt.Printf("%d/%d: call %s\n",currentUrlIndex, len(assets)-1, assetUrl)

		go fetchAsset(assetUrl, client, chUrls, chFinished, &o.verbose)

		//limit calls count to max_concurrent_call
		currentUrlIndex++
		if currentUrlIndex == *parallelFetch {
			break
		}
	}

	// Subscribe to channels to wait for go routine
	for c := 0; c < len(assets); {
		select {
		case stat := <-chUrls:
			assetsStats = append(assetsStats, stat)
			gstat.totalResponseSize += stat.responseSize
		//got an asset, fetch next if exist
		case <-chFinished:
			if currentUrlIndex < len(assets) {
				go fetchAsset(assets[currentUrlIndex], client, chUrls, chFinished, &o.verbose)
			}
			c++
			currentUrlIndex++
		}
	}

	close(chUrls)

	//Set timer for global time
	gstat.totalResponseTime += time.Now().Sub(t0)

	return gstat
}
