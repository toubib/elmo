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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"github.com/fatih/color"
)

type downloadStatistic struct {
	url           string
	responseTime time.Duration
	responseSize int
	statusCode   int
}

type globalStatistic struct {
	totalResponseTime time.Duration
	totalResponseSize int
}

var (
	VERSION    = "0.1-dev"
	BUILD_DATE = ""
)

//colors !
var (
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan = color.New(color.FgCyan).SprintFunc()
	white = color.New(color.FgWhite).SprintFunc()
)

//cli flags
var (
	url = flag.String("url", "", "The url to get.")
	version = flag.Bool("version", false, "Print version information.")
	verbose = flag.Bool("verbose", false, "Print more informations.")
	parallelFetch = flag.Int("parallel", 8, "Number of parallel fetch to launch. 0 means unlimited.")
)

var globalStartTime = time.Now()

// Helper function to pull the  attribute from a Token
func getLink(t html.Token) (ok bool, link string) {

	// Check link types, we need only stylesheet
	if t.Data == "link" {
		for _, a := range t.Attr {
			if a.Key == "rel" && a.Val != "stylesheet" {
				ok = false
				return
			}
		}

	}

	// Iterate over all of the Token's attributes until we find an "src"
	for _, a := range t.Attr {
		if a.Key == "src" || a.Key == "href" {
			link = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Extract all http** links from a given webpage
func fetchMainUrl(url string) ([]string, downloadStatistic) {

	//List of urls found
	var assets []string

	//set downloadStatistic
	stat := downloadStatistic{url, 0, 0, 0}

	//timer before
	t0 := time.Now()

	//launch the query
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to get input url \"" + url + "\"")
		return assets, stat
	}

	//Set stats
	stat.responseTime = time.Now().Sub(t0)
	stat.statusCode = resp.StatusCode

	//get the body size
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ERROR: Failed to read body for input url \"" + url + "\"")
		return assets, stat
	}

	//Set response size stat
	stat.responseSize = len(body)

	//Print download
	if (*verbose) {
		fmt.Printf("%s\t%s %s %v %v%s\n", time.Now().Sub(globalStartTime), yellow(stat.statusCode), stat.url, cyan(stat.responseTime), white(stat.responseSize),white("b"))
	}

	//extract assets from html
	assets = extractAssets(body)

	return assets, stat
}


//Get a html body and extract all assets links
func extractAssets(body []byte) []string {
	var assets []string

	//create the tokenizer
	z := html.NewTokenizer(bytes.NewReader(body))

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			//fmt.Printf("   - end of doc\n")
			return assets
		case html.SelfClosingTagToken, html.StartTagToken:
			t := z.Token()

			// Check if the token is a target tag
			// And extract the link if there is one
			if t.Data != "img" &&
				t.Data != "script" &&
				t.Data != "embed" &&
				t.Data != "link" && //only stylesheet
				t.Data != "input" {

				continue
			}

			linkFound, url := getLink(t)

			// We do not found a link
			if !linkFound {
				continue
			}

			// Make sure the url begines in http**
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				assets = append(assets, url)
			}
		}
	}
}

//Fetch an asset and get downloadStatistic
func fetchAsset(url string, chStat chan downloadStatistic, chFinished chan bool) {

	//set downloadStatistic
	stat := downloadStatistic{url, 0, 0, 0}

	//timer before
	t0 := time.Now()

	//launch the query
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to get link \"" + url + "\"")
		return
	}

	//Set stat
	stat.responseTime = time.Now().Sub(t0)
	stat.statusCode = resp.StatusCode

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	//get the body size
	b := resp.Body
	defer b.Close() // close Body when the function returns
	body, err := ioutil.ReadAll(resp.Body)

	//Set response size stat
	stat.responseSize = len(body)

	//Print download
	if (*verbose) {
		fmt.Printf("%s\t%s %s %v %v%s\n", time.Now().Sub(globalStartTime), yellow(stat.statusCode), stat.url, cyan(stat.responseTime), white(stat.responseSize),white("b"))
	}

	chStat <- stat
}

func main() {
	//urls and global stats
	var assets []string
	var assetsStats []downloadStatistic
	var mainUrlStat downloadStatistic
	var gstat globalStatistic
	var currentUrlIndex int

	flag.Parse()

	// Channels
	chUrls := make(chan downloadStatistic)
	chFinished := make(chan bool)

	// Manage flags stuff
	if *version {
		fmt.Printf("%v\nBuild: %v\n", VERSION, BUILD_DATE)
		return
	}

	//Set timer for global time
	t0 := time.Now()

	//Fetch the main url and get inner links
	assets, mainUrlStat = fetchMainUrl(*url)
	gstat.totalResponseSize += mainUrlStat.responseSize

	//Fetch the firsts inner links
	for _, url := range assets {
		//fmt.Printf("%d/%d: call %s\n",currentUrlIndex, len(assets)-1, url)

		go fetchAsset(url, chUrls, chFinished)

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
			assetsStats = append(assetsStats,stat)
			gstat.totalResponseSize += stat.responseSize
		//got an asset, fetch next if exist
		case <-chFinished:
			if currentUrlIndex < len(assets) {
				go fetchAsset(assets[currentUrlIndex], chUrls, chFinished)
			}
			c++
			currentUrlIndex++
		}
	}

	close(chUrls)

	//Set timer for global time
	gstat.totalResponseTime += time.Now().Sub(t0)

	//Use that to send data to influx
	//for _, stat := range assetsStats {
	//	if (*verbose) {
	//		fmt.Printf("%s\t%s %s %v %v%s\n", time.Now().Sub(globalStartTime), yellow(stat.statusCode), stat.url, cyan(stat.responseTime), white(stat.responseSize),white("b"))
	//	}

	//}

	// We're done! Print the results...
	fmt.Printf("The call took %v to run.\n", cyan(gstat.totalResponseTime))
	fmt.Printf("Total size: %v%s.\n", white(gstat.totalResponseSize/1024), white("kb"))

}
