package main

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/json-iterator/go"
	"github.com/tkanos/gonfig"
	"log"
	"strconv"
	"time"
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

type Configuration struct {
	Host     string
	Database string
	Username string
	Password string
	Url      string
}

func main() {

	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		fmt.Println("Error gonfig: ", err)
		log.Fatal(err)
	}

	_host := configuration.Host
	_db := configuration.Database
	_username := configuration.Username
	_password := configuration.Password
	_url := configuration.Url

	for {
		fmt.Println(time.Now())

		// GET request
		resp, err := resty.R().Get(_url)
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
				Addr:     _host,
				Username: _username,
				Password: _password,
			})
			if err != nil {
				log.Fatal(err)
			}

			// Create a new point batch
			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  _db,
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
