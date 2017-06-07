package main

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/json-iterator/go"
	"log"
	"strconv"
	"time"
	"os"
)

type Station struct {
	ID             string `json:"id"`
	District       string `json:"district"`
	Lon            string `json:"lon"`
	Lat            string `json:"lat"`
	Bikes          string `json:"bikes"`
	Slots          string `json:"slots"`
	Zip            string `json:"zip"`
	Address        string `json:"address"`
	AddressNumber  string `json:"addressNumber"`
	NearbyStations string `json:"nearbyStations"`
	Status         string `json:"status"`
	Name           string `json:"name"`
}

func main() {

	h := os.Getenv("INFLUXDB_HOST")
	fmt.Printf("INFLUXDB_HOST: %s\n", h)

	d := os.Getenv("INFLUXDB_DATABASE")
	fmt.Printf("INFLUXDB_DATABASE: %s\n", d)

	u := os.Getenv("INFLUXDB_USERNAME")
	fmt.Printf("INFLUXDB_USERNAME: %s\n", u)

	p := os.Getenv("INFLUXDB_PASSWORD")
	fmt.Printf("INFLUXDB_PASSWORD: %s\n", p)

	v := os.Getenv("VELO_URL")
	fmt.Printf("VELO_URL: %s\n", v)

	for {
		fmt.Println(time.Now())

		// GET request
		resp, err := resty.R().Get(v)
		if err != nil {
			fmt.Println("Error GET: ", err)
			return
		}

		var Stations []Station

		JsonObject := resp.Body()
		JsonBytes := []byte(JsonObject)

		err = jsoniter.Unmarshal(JsonBytes, &Stations)
		if err != nil {
			fmt.Println("Error jsoniter: ", err)
		}

		for idx := range Stations {
			// Create a new HTTPClient
			c, err := client.NewHTTPClient(client.HTTPConfig{
				Addr:     h,
				Username: u,
				Password: p,
			})
			if err != nil {
				log.Fatal(err)
			}

			// Create a new point batch
			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  d,
				Precision: "s",
			})
			if err != nil {
				fmt.Println("Error NewBatchPoints: ", err)
				log.Fatal(err)
			}

			IBikes, err := strconv.Atoi(Stations[idx].Bikes)
			if err != nil {
				fmt.Println("Error strconv bikes: ", err)
				log.Fatal(err)
			}

			ISlots, err := strconv.Atoi(Stations[idx].Slots)
			if err != nil {
				fmt.Println("Error strconv slots: ", err)
				log.Fatal(err)
			}

			Total := IBikes + ISlots

			// Create a point and add to batch
			tags := map[string]string{"Station": Stations[idx].Name}
			fields := map[string]interface{}{
				"Bikes": IBikes,
				"Slots": ISlots,
				"Total": int(Total),
			}

			pt, err := client.NewPoint("Velo", tags, fields, time.Now())
			if err != nil {
				fmt.Println("Error NewPoint: ", err)
				log.Fatal(err)
			}
			bp.AddPoint(pt)

			// Write the batch
			if err := c.Write(bp); err != nil {
				fmt.Println("Error c.write: ", err)
				log.Fatal(err)
			}
			c.Close()
		}
		time.Sleep(10 * time.Second)
	}
}
