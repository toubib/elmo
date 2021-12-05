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
	"context"
	"errors"

	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/html"
)

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

const (
	NAGIOS_OK      = 0
	NAGIOS_WARNING = 1
	NAGIOS_ERROR   = 2
	NAGIOS_UNKNOWN = 3
)

var (
	VERSION = "0.3-dev"
)

//colors !
var (
	//yellow     = color.New(color.FgYellow).SprintFunc()
	cyan       = color.New(color.FgCyan).SprintFunc()
	white      = color.New(color.FgWhite).SprintFunc()
	bold_white = color.New(color.FgWhite, color.Bold).SprintFunc()
	red        = color.New(color.FgRed, color.Bold).SprintFunc()
	green      = color.New(color.FgGreen, color.Bold).SprintFunc()
)

//cli flags
var (
	debug     = false
	verbose   = false
	useNagios bool
	timeout   int
)

func cliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Usage:    "The url to get",
			Aliases:  []string{"u"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "user-agent",
			Usage:   "Change the user-agent",
			Aliases: []string{"A"},
		},
		&cli.StringFlag{
			Name:    "keyword",
			Usage:   "Check for keyword in reponse",
			Aliases: []string{"K"},
		},
		&cli.BoolFlag{
			Name:        "debug",
			Value:       false,
			Destination: &debug,
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Value:       false,
			Destination: &verbose,
		},
		&cli.IntFlag{
			Name:    "parallel",
			Value:   8,
			Usage:   "Number of parallel fetch to launch. 0 means unlimited",
			Aliases: []string{"p"},
		},
		&cli.IntFlag{
			Name:  "connect-timeout",
			Value: 1000,
			Usage: "Connect timeout in ms",
		},
		&cli.IntFlag{
			Name:  "tls-timeout",
			Value: 1000,
			Usage: "TLS handshake timeout in ms",
		},
		&cli.StringFlag{
			Name:  "resolve",
			Usage: "<host:port:addr> Resolve the host+port to this address",
		},
		&cli.BoolFlag{
			Name:        "use-nagios",
			Value:       false,
			Usage:       "Nagios compatible output.",
			Destination: &useNagios,
		},
		&cli.IntFlag{
			Name:  "nagios-warning",
			Value: 5000,
			Usage: "Nagios warning time in ms",
		},
		&cli.IntFlag{
			Name:  "nagios-critical",
			Value: 10000,
			Usage: "Nagios critical time in ms",
		},
		&cli.IntFlag{
			Name:        "timeout",
			Value:       10000,
			Usage:       "Global request timeout in ms.",
			Destination: &timeout,
			Aliases:     []string{"t"},
		},
		&cli.IntFlag{
			Name:  "response-header-timeout",
			Value: 0,
			Usage: "Response header timeout in ms",
		},
		&cli.BoolFlag{
			Name:  "use-influx",
			Value: false,
			Usage: "Send data to influxdb",
		},
		&cli.StringFlag{
			Name:  "influx-url",
			Usage: "The influx database access url",
			Value: "http://localhost:8086",
		},
		&cli.StringFlag{
			Name:  "influx-database",
			Usage: "The influx database name",
			Value: "elmo",
		},
		&cli.StringFlag{
			Name:  "assets-allowed-domains",
			Usage: "List of allowed assets domains to fetch from, comma separated",
		},
	}
}

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
func fetchMainUrl(mainUrl string, client *http.Client, headers map[string]string, keyword string) ([]string, downloadStatistic, error) {

	//List of urls found
	var assets []string

	//set downloadStatistic
	stat := downloadStatistic{mainUrl, 0, 0, 0}

	//timer before
	t0 := time.Now()

	//launch the query
	req, _ := http.NewRequest("GET", mainUrl, nil)

	//set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if debug {
		fmt.Printf("debug request: %v\n", req)
	}

	if debug || verbose {

		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			fmt.Printf("Redirect to %v\n", req.URL)
			return nil
		}
	}

	resp, err := client.Do(req)

	if debug {
		fmt.Printf("debug response: %v\n", resp)
	}

	if err != nil {
		return assets, stat, err
	}

	//Set stats
	stat.responseTime = time.Since(t0)
	stat.statusCode = resp.StatusCode

	//get the body size
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return assets, stat, err
	}

	//Check for keyword
	if keyword != "" && !bytes.Contains(body, []byte(keyword)) {
		return assets, stat, errors.New("String " + keyword + " not found.")
	}

	//Set response size stat
	stat.responseSize = len(body)

	//Print download
	if verbose {
		fmt.Printf("%s\t%s %s %v %v%s\n", time.Since(globalStartTime), green(stat.statusCode), stat.url, cyan(stat.responseTime), white(stat.responseSize), white("b"))
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
func fetchAsset(assetUrl string, assetsAllowedDomains string, client *http.Client, headers map[string]string, chStat chan downloadStatistic, chFinished chan bool) {

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

	if !checkIfDomainAllowed(assetsAllowedDomains, &req.URL.Host) {
		return
	}

	//set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	//handle error
	if err != nil {
		if !useNagios {
			fmt.Println(red("Error:"), stat.url, err)
		}
		return
	}

	//Set stat
	stat.responseTime = time.Since(t0)
	stat.statusCode = resp.StatusCode

	//get the body size
	b := resp.Body
	defer b.Close() // close Body when the function returns
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if !useNagios {
			fmt.Println(red("Error:"), stat.url, err)
		}
		stat.responseSize = 0
	} else {
		//Set response size stat
		stat.responseSize = len(body)
	}

	//Print download
	if verbose {
		fmt.Printf("%s\t%s %s %v %v%s\n", time.Since(globalStartTime), green(stat.statusCode), stat.url, cyan(stat.responseTime), white(stat.responseSize), white("b"))
	}

	chStat <- stat
}

// test if the given domain is allowed to fetch
func checkIfDomainAllowed(assetsAllowedDomains string, host *string) bool {

	if assetsAllowedDomains == "" {
		return true
	}

	allowedDomains := strings.Split(assetsAllowedDomains, ",")

	for _, domain := range allowedDomains {
		if domain == *host {
			return true
		}
	}

	return false
}

func main() {

	app := cli.NewApp()
	app.Name = "elmo"
	app.Usage = "Elmo web client"
	app.Version = VERSION
	app.Flags = cliFlags()

	app.Action = func(cli *cli.Context) error {

		//urls and global stats
		var (
			assets          []string
			assetsStats     []downloadStatistic
			mainUrlStat     downloadStatistic
			gstat           globalStatistic
			currentUrlIndex int
			err             error
		)

		// Channels
		chUrls := make(chan downloadStatistic)
		chFinished := make(chan bool)

		//set global timeout if nagios timeout is set
		if cli.Bool("use-nagios") {
			timeout = cli.Int("nagios-critical")
		}

		//Set headers
		headers := make(map[string]string)
		if cli.String("user-agent") != "" {
			headers["User-Agent"] = cli.String("user-agent")
		}

		//set timeouts
		transport := &http.Transport{

			DialContext: (&net.Dialer{
				Timeout: time.Duration(cli.Int("connect-timeout")) * time.Millisecond,
			}).DialContext,

			TLSHandshakeTimeout:   time.Duration(cli.Int("tls-timeout")) * time.Millisecond,
			ResponseHeaderTimeout: time.Duration(cli.Int("response-header-timeout")) * time.Millisecond,
		}

		dialer := &net.Dialer{
			Timeout: time.Duration(cli.Int("connect-timeout")) * time.Millisecond,
			//        KeepAlive: 30 * time.Second,
		}

		var domain_resolve []string = nil

		if cli.String("resolve") != "" {
			domain_resolve = strings.Split(cli.String("resolve"), ":")
			if len(domain_resolve) != 3 {
				fmt.Printf("bad argument -resolve\n")
				os.Exit(NAGIOS_UNKNOWN)
			}
			if debug {
				fmt.Printf("debug: domain_resolve set to %v\n", domain_resolve)
			}
		}

		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			if cli.String("resolve") != "" {
				if addr == domain_resolve[0]+":"+domain_resolve[1] {
					if debug {
						fmt.Printf(bold_white("debug: rewrite %s to %s\n"),
							addr, domain_resolve[2]+":"+domain_resolve[1])
					}
					addr = domain_resolve[2] + ":" + domain_resolve[1]
				}
			}
			return dialer.DialContext(ctx, network, addr)
		}

		//Set an http client with this transport
		client := &http.Client{
			Timeout:   time.Duration(timeout) * time.Millisecond,
			Transport: transport,
		}

		//Set timer for global time
		t0 := time.Now()

		//Fetch the main url and get inner links
		assets, mainUrlStat, err = fetchMainUrl(cli.String("url"), client, headers, cli.String("keyword"))
		assetsStats = append(assetsStats, mainUrlStat)

		//handle main url error
		if err != nil {
			if cli.Bool("use-nagios") {
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

			go fetchAsset(assetUrl, cli.String("assets-allowed-domains"),
				client, headers, chUrls, chFinished)

			//limit calls count to max_concurrent_call
			currentUrlIndex++
			if currentUrlIndex == cli.Int("parallel") {
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
					go fetchAsset(assets[currentUrlIndex], cli.String("assets-allowed-domains"),
						client, headers, chUrls, chFinished)
				}
				c++
				currentUrlIndex++
			}
		}

		close(chUrls)

		//Set timer for global time
		gstat.totalResponseTime += time.Since(t0)

		// send data to influxdb
		if cli.Bool("use-influx") {
			sendstatsToInflux(cli.String("influx-url"), cli.String("influx-database"),
				cli.String("url"), &assetsStats)
		}

		// We're done! Print the results...
		if cli.Bool("use-nagios") {
			fmt.Printf("Downloaded %vKB in %d/%d files in %v.|size=%vKB time=%v;%v;%v;0;%v\n",
				gstat.totalResponseSize/1024, len(assetsStats), len(assets), gstat.totalResponseTime,
				gstat.totalResponseSize/1024, gstat.totalResponseTime,
				cli.Int("nagios-warning"), cli.Int("nagios-critical"), cli.Int("timeout"),
			)

			//nagios exit
			if gstat.totalResponseTime >= time.Duration(cli.Int("nagios-critical"))*time.Millisecond {
				os.Exit(NAGIOS_ERROR)
			} else if gstat.totalResponseTime >= time.Duration(cli.Int("nagios-warning"))*time.Millisecond {
				os.Exit(NAGIOS_WARNING)
			} else {
				os.Exit(NAGIOS_OK)
			}

		} else {
			fmt.Printf("Downloaded assets: %d/%d.\n", len(assetsStats), len(assets))
			fmt.Printf("Total time: %v.\n", cyan(gstat.totalResponseTime))
			fmt.Printf("Total size: %v%s.\n", white(gstat.totalResponseSize/1024), white("kb"))
		}

		return nil
	}

	app.Run(os.Args)

}
