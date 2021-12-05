package main

import (
	"fmt"
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

//Send statistic data to influxdb
func sendstatsToInflux(influxUrl string, influxDatabase string, mainUrl string, assetsStats *[]downloadStatistic) {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: influxUrl})
	if err != nil {
		fmt.Println(red("Influxdb - error creating InfluxDB Client:\n"), err.Error())
	}

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDatabase,
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
		pt, err := client.NewPoint(mainUrl, tags, fields, influxTime)
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
