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

type download_statistic struct {
	url           string
	response_time time.Duration
	response_size int
}

type global_statistic struct {
	total_response_time time.Duration
	total_response_size int
}

var (
	VERSION    = "0.1-dev"
	BUILD_DATE = ""
)

var url = flag.String("url", "", "The url to get")
var version = flag.Bool("version", false, "print version information")

// Helper function to pull the  attribute from a Token
func getSrc(t html.Token) (ok bool, src string) {
	// Iterate over all of the Token's attributes until we find an "src"
	for _, a := range t.Attr {
		if a.Key == "src" {
			src = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Extract all http** links from a given webpage
func fetch_main_url(url string) (map[string]bool, download_statistic) {

	//List of urls found
	foundUrls := make(map[string]bool)

	//set download_statistic
	stat := download_statistic{url, 0, 0}

	//timer before
	t0 := time.Now()

	//launch the query
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to get input url \"" + url + "\"")
		return foundUrls, stat
	}

	//timer after
	t1 := time.Now()

	//Set request time stat
	stat.response_time = t1.Sub(t0)

	//get the body size
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ERROR: Failed to read body for input url \"" + url + "\"")
		return foundUrls, stat
	}

	//Set response size stat
	stat.response_size = len(body)

	//Print download
	fmt.Printf(" - [%s] %s %v %v\n", resp.Status, stat.url, stat.response_time, stat.response_size)

	//create the tokenizer
	z := html.NewTokenizer(bytes.NewReader(body))

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			//fmt.Printf("   - end of doc\n")
			return foundUrls, stat
		case tt == html.SelfClosingTagToken:
			t := z.Token()

			// Check if the token is an <img> tag
			isAnchor := t.Data == "img"
			if !isAnchor {
				continue
			}

			// Extract the src value, if there is one
			ok, url := getSrc(t)
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				foundUrls[url] = true
			}
		}
	}
}

//Fetch an asset and get download_statistic
func fetch_asset(url string, chStat chan download_statistic, chFinished chan bool) {

	//set download_statistic
	stat := download_statistic{url, 0, 0}

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
	stat.response_time = t1.Sub(t0)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	//get the body size
	b := resp.Body
	defer b.Close() // close Body when the function returns
	body, err := ioutil.ReadAll(resp.Body)

	//Set response size stat
	stat.response_size = len(body)

	//Print download
	fmt.Printf(" - [%s] %s %v %v\n", resp.Status, stat.url, stat.response_time, stat.response_size)

	chStat <- stat
}

func main() {
	//urls and global stats
	var foundUrls map[string]bool
	var main_url_stat download_statistic
	var gstat global_statistic

	flag.Parse()

	// Channels
	chUrls := make(chan download_statistic)
	chFinished := make(chan bool)

	// Manage flags stuff
	if *version {
		fmt.Printf("%v\nBuild: %v\n", VERSION, BUILD_DATE)
		return
	}

	//Set timer for global time
	t0 := time.Now()

	//Fetch the main url and get inner links
	foundUrls, main_url_stat = fetch_main_url(*url)
	gstat.total_response_time += main_url_stat.response_time
	gstat.total_response_size += main_url_stat.response_size

	//Fetch inner links
	for url, _ := range foundUrls {
		go fetch_asset(url, chUrls, chFinished)
	}

	// Subscribe to channels to wait for go routine
	for c := 0; c < len(foundUrls); {
		select {
		case stat := <-chUrls:
			gstat.total_response_time += stat.response_time
			gstat.total_response_size += stat.response_size
			foundUrls[stat.url] = true
		case <-chFinished:
			c++
		}
	}

	//Set timer for global time
	t1 := time.Now()

	// We're done! Print the results...
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
	fmt.Printf("Cumulated time: %v.\n", gstat.total_response_time)
	fmt.Printf("Cumulated size: %v.\n", gstat.total_response_size)

	close(chUrls)
}
