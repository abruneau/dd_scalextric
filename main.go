package main

import (
	"flag"
	"log"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/abruneau/dd_scalextric/race"
	"github.com/abruneau/dd_scalextric/utils"
)

var configPath = flag.String("config", "./config.yaml", "the config file with gpio ports")

func main() {
	c, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("StatsD client started")

	var config utils.Configuration

	err = config.Get(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	race.InitRace(&config, c)
}
