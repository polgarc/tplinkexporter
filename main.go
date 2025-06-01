package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/thelastguardian/tplinkexporter/clients"
	"github.com/thelastguardian/tplinkexporter/collectors"
)

func main() {
	var (
		host     = kingpin.Flag("host", "Host of target tplink easysmart switch.").Envar("HOST").Required().String()
		username = kingpin.Flag("username", "Username for switch GUI login").Envar("USERNAME").Default("admin").String()
		password = kingpin.Flag("password", "Password for switch GUI login").Envar("PASSWORD").Required().String()
		port     = kingpin.Flag("port", "Metrics port to listen on for prometheus scrapes").Envar("PORT").Default("9717").Int()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	tplinkSwitch, err := clients.NewTPLinkSwitch(*host, *username, *password)
	if err != nil {
		log.Fatalf("Error creating tplink switch client: %v", err)
	}
	trafficCollector := collectors.NewTrafficCollector("tplinkexporter", tplinkSwitch)
	prometheus.MustRegister(trafficCollector)
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Beginning to serve on port %v", *port)
	err = tplinkSwitch.Login()
	if err != nil {
		log.Println(err)
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
