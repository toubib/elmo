//http://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html

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
)

type downloadStatistic struct {
	url           string
	responseTime time.Duration
	responseSize int
}

type globalStatistic struct {
	totalResponseTime time.Duration
	totalResponseSize int
}

var (
	VERSION    = "0.1-dev"
	BUILD_DATE = ""
)

var url = flag.String("url", "", "The url to get.")
var version = flag.Bool("version", false, "Print version information.")
var parallelFetch = flag.Int("parallel", 8, "Number of parallel fetch to launch. 0 means unlimited.")

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
	stat := downloadStatistic{url, 0, 0}

	//timer before
	t0 := time.Now()

	//launch the query
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to get input url \"" + url + "\"")
		return assets, stat
	}

	//timer after
	t1 := time.Now()

	//Set request time stat
	stat.responseTime = t1.Sub(t0)

	//get the body size
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ERROR: Failed to read body for input url \"" + url + "\"")
		return assets, stat
	}

	//Set response size stat
	stat.responseSize = len(body)

	//Print download
	fmt.Printf(" - [%s] %s %v %v\n", resp.Status, stat.url, stat.responseTime, stat.responseSize)

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
	stat := downloadStatistic{url, 0, 0}

	//timer before
	t0 := time.Now()

	//launch the query
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to get link \"" + url + "\"")
		return
	}

	//timer after
	t1 := time.Now()

	//Set request time stat
	stat.responseTime = t1.Sub(t0)

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
	fmt.Printf(" - [%s] %s %v %v\n", resp.Status, stat.url, stat.responseTime, stat.responseSize)

	chStat <- stat
}

func main() {
	//urls and global stats
	var assets []string
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
	gstat.totalResponseTime += mainUrlStat.responseTime
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
			gstat.totalResponseTime += stat.responseTime
			gstat.totalResponseSize += stat.responseSize
		//got an asset, fetch next if exist
		case <-chFinished:
			if currentUrlIndex < len(assets) {
				//fmt.Printf("%d/%d: call %s\n",currentUrlIndex, len(assets)-1, assets[currentUrlIndex])
				go fetchAsset(assets[currentUrlIndex], chUrls, chFinished)
			}
			c++
			currentUrlIndex++
		}
	}

	//Set timer for global time
	t1 := time.Now()

	// We're done! Print the results...
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
	fmt.Printf("Cumulated time: %v.\n", gstat.totalResponseTime)
	fmt.Printf("Cumulated size: %v.\n", gstat.totalResponseSize)

	close(chUrls)
}
