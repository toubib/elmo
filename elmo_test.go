package main

import (
	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
)

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
	}{ {5, 385} }

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, htmlBody)
	}))
	defer ts.Close()

	for _, tt := range tests {
		assets, mainUrlStat := fetchMainUrl(ts.URL)

		if (len(assets) != tt.assetCount){
			t.Errorf("fetchMainUrl do not returned %d elements but %d", tt.assetCount,len(assets))
		}

		if (mainUrlStat.responseSize != tt.responseSize){
			t.Errorf("mainUrlStat.responseSize is not returned %d but %d",tt.responseSize,mainUrlStat.responseSize)
		}
	}
}
