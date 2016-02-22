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
	"fmt"
	"github.com/mreiferson/go-httpclient"
	"net/http"
	"net/http/httptest"
//	"net/url"
	"testing"
	"time"
)

//TODO add test dl 3 differents images and check images sizes and total size
//bug: all assets are sames :(

func TestFetchMainUrl(t *testing.T) {

	const htmlBody = `<head>
		<link rel="alternate" type="application/rss+xml" title="Flux" href="http://test.com/feed/" />
		<link rel="stylesheet" href="http://test.com/1.css" type="text/css" media="all" />
		</head><body>
		<img src="http://test.com/1.png"/>
		<img src="http://test.com/É.png"/>
		<img src="http://test.com/3.png">
		<script type="text/javascript" src="http://test.com/1.js">
	</body>`

	tests := []struct {
		assetCount, responseSize int
	}{{5, 385}}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/":
			fmt.Fprintln(w, htmlBody)
		}
	}))
	defer ts.Close()

	//Set a transport with timeouts
	transport := &httpclient.Transport{
		ConnectTimeout:        1 * time.Second,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	defer transport.Close()

	//Set an http client with this transport
	client := &http.Client{Transport: transport}

	for _, tt := range tests {
		assets, mainUrlStat := fetchMainUrl(&ts.URL, client)

		if len(assets) != tt.assetCount {
			t.Errorf("fetchMainUrl do not returned %d elements but %d", tt.assetCount, len(assets))
		}

		if mainUrlStat.responseSize != tt.responseSize {
			t.Errorf("mainUrlStat.responseSize is not returned %d but %d", tt.responseSize, mainUrlStat.responseSize)
		}
	}
}

func TestFetchAsset(t *testing.T) {
	tests := []struct {
		assetUrl string
		responseSize int
	}{
		{"/1.png", 1},
		{"/É.png", 2},
		{"/3.png", 3},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/1.png":
			fmt.Fprintln(w, "0")
		case "/É.png":
			fmt.Fprintln(w, "0000")
		case "/3.png":
			fmt.Fprintln(w, "00000000")
		}
	}))
	defer ts.Close()

	//Set a transport with timeouts
	transport := &httpclient.Transport{
		ConnectTimeout:        1 * time.Second,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	defer transport.Close()

	//Set an http client with this transport
	client := &http.Client{Transport: transport}

	// Channels
	chUrls := make(chan downloadStatistic)
	chFinished := make(chan bool)

	for _, tt := range tests {

		//make absolute url
		u := ts.URL + tt.assetUrl

		// fetch asset
		go fetchAsset(u, client, chUrls, chFinished)
	}

	// Subscribe to channels to wait for go routine
	for c := 0; c < len(tests); {
		select {
		case stat := <-chUrls:
			fmt.Println(tests[0])
			fmt.Println("stat: ",stat)
		//	assetsStats = append(assetsStats, stat)
		//	gstat.totalResponseSize += stat.responseSize
		//got an asset, fetch next if exist
		case <-chFinished:
			c++
		}
	}

	close(chUrls)

	//if mainUrlStat.responseSize != tt.responseSize {
	//	t.Errorf("mainUrlStat.responseSize is not returned %d but %d", tt.responseSize, mainUrlStat.responseSize)
	//}

}

func TestCheckIfDomainAllowed(t *testing.T) {

	tests := []struct {
		url    string
		result bool
	}{
		{"http://test.com", true},
		{"https://test.com", true},
		{"http://test2.com", false},
		{"http://test3.com", true},
		{"http://test4.com", false},
	}

	//force assets-allowed-domains flag
	*assetsAllowedDomains = "test.com,test3.com"

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.url, nil)

		testResult := (checkIfDomainAllowed(&req.URL.Host) == tt.result)
		if !testResult {
			t.Errorf("checkIfDomainAllowed (%v) has not returned %v but %v", tt.url, tt.result, testResult)
		}
	}
}
