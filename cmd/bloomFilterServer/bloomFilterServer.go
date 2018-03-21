/* This file will contain functionalities to host the main service
in which the client will make requests to. Use net/http to create
API endpoints to access the following functionalities:
	- UpdateBloomFilter
		- AddToBloomFilter
		- RepopulateBloomFilter
	- GetArrayOfUnsubscribedEmails

	- when the http endpoint is hit go to the graphana to graph those metrics
	- use docker to integrate with AWS?
	- graphana and graphite have docker images
	- mysql database is currently local buttttt, should spin up a socker container with mysql
		- so we're all on the same page
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/spf13/viper"
	"github.com/vlam321/Inf191BloomFilter/bloomManager"
	"github.com/vlam321/Inf191BloomFilter/databaseAccessObj"
	"github.com/vlam321/Inf191BloomFilter/payload"

	metrics "github.com/rcrowley/go-metrics"
)

//The bloom filter for this server
var bf *bloomManager.BloomFilter
var shard int
var host string

// handleDetermineUpdate
func handleDetermineUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %v %v %v\n", r.Method, r.URL, r.Proto)
	bf.DetermineUpdateOrRepopulate(shard)
}

// handleFilterUnsubscribed
func handleFilterUnsubscribed(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %v %v %v\n", r.Method, r.URL, r.Proto)
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: Unable to read request data. %v\n", err)
		return
	}

	var pl payload.Payload
	err = json.Unmarshal(bytes, &pl)
	if err != nil {
		log.Printf("Error: Unable to unmarshal Payload. %v\n", err)
		return
	}

	//uses bloomManager to get the result of unsubscribed emails
	//puts them in struct, result
	filteredResults := bf.GetArrayOfUnsubscribedEmails(map[int][]string{pl.UserId: pl.Emails})
	results := payload.Payload{pl.UserId, filteredResults[pl.UserId]}

	jsn, err := json.Marshal(results)
	if err != nil {
		log.Printf("Error marshaling filtered emails. %v\n", err)
		return
	}

	metrics.GetOrRegisterCounter(fmt.Sprintf("%s.request.numreq", host), nil).Inc(1)
	//write back to client
	w.Write(jsn)
}

// handleQueryUnsubscribed
func handleQueryUnsubscribed(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %v %v %v\n", r.Method, r.URL, r.Proto)
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: Unable to read request data. %v\n", err)
		return
	}

	var pl payload.Payload
	err = json.Unmarshal(bytes, &pl)
	if err != nil {
		log.Printf("Error: Unable to unmarshal Payload. %v\n", err)
		return
	}

	//uses bloomManager to get the result of unsubscribed emails
	//puts them in struct, result
	filteredResults := bf.QueryUnsubscribed(map[int][]string{pl.UserId: pl.Emails})
	results := payload.Payload{pl.UserId, filteredResults[pl.UserId]}

	jsn, err := json.Marshal(results)
	if err != nil {
		log.Printf("Error marshaling filtered emails. %v\n", err)
		return
	}

	metrics.GetOrRegisterCounter(fmt.Sprintf("%s.request.numreq", host), nil).Inc(1)
	//write back to client
	w.Write(jsn)
}

//updateBloomFilterBackground manages the periodic updates of the bloom
//filter. Update calls repopulate, creating a new updated bloom filter
func updateBloomFilterBackground(dao *databaseAccessObj.Conn) {
	//Set new ticker to every 2 seconds
	ticker := time.NewTicker(time.Second * 3)

	for t := range ticker.C {
		//Call update bloom filter
		metrics.GetOrRegisterGauge(fmt.Sprintf("%s.dbsize.gauge", host), nil).Update(int64(dao.GetTableSize(shard)))
		bf.RepopulateBloomFilter(shard)
		log.Println("Bloom Filter Updated at: ", t.Format("2006-01-02 3:4:5 PM"))
	}
}

func graphDBSize(dao *databaseAccessObj.Conn) {
	//Set new ticker to every 2 seconds
	ticker := time.NewTicker(time.Second * 60)

	for _ = range ticker.C {
		//Call update bloom filter
		metrics.GetOrRegisterGauge(fmt.Sprintf(fmt.Sprintf("%s.dbsize.gauge", host)), nil).Update(int64(dao.GetTableSize(shard)))
	}
}

// setBloomFilter initialize bloom filter
func setBloomFilter(dao *databaseAccessObj.Conn) {
	numEmails := uint(dao.GetTableSize(shard))
	falsePositiveProb := float64(0.001)
	bf = bloomManager.New(numEmails, falsePositiveProb)
}

// Retrieve the IPv4 address of the current AWS EC2 instance
func getMyIP() (string, error) {
	resp, err := http.Get("http://checkip.amazonaws.com/")
	if err != nil {
		return "x.x.x.x", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "x.x.x.x", err
	}
	return string(body[:]), nil
}

func mapBf2Shard() error {
	viper.SetConfigName("bfIPConf")
	viper.AddConfigPath("settings")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := mapBf2Shard()
	if err != nil {
		log.Printf("BloomFilter: %v\n", err)
		return
	}

	if viper.GetString("host") == "docker" {
		tabnum, err := strconv.Atoi(os.Getenv("SHARD"))
		if err != nil {
			log.Printf("Bloom Filter: %v\n", err)
		}
		shard = tabnum
	} else if viper.GetString("host") == "ecs" {
		shard, err = strconv.Atoi(os.Getenv("SHARD"))
		if err != nil {
			log.Printf("Bloom Filter: %v\n", err)
		}
	} else {
		log.Printf("BloomFilter: Invalid host config.")
		return
	}

	log.Printf("SUCCESSFULLY MAPPED FILTER TO DB SHARD.\n")
	log.Printf("HOSTING ON: %s\n", viper.GetString("host"))
	log.Printf("USING SHARD: %d\n", shard)

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:2003", os.Getenv("GRAPHITE_IP")))
	host, _ = os.Hostname()
	go graphite.Graphite(metrics.DefaultRegistry, 10e9, "metrics", addr)

	dao := databaseAccessObj.New()
	defer dao.CloseConnection()

	setBloomFilter(dao)
	bf.RepopulateBloomFilter(shard)
	//Run go routine to make periodic updates
	//Runs until the server is stopped
	//go updateBloomFilterBackground(dao)
	go graphDBSize(dao)
	http.HandleFunc("/update", handleDetermineUpdate)
	http.HandleFunc("/filterUnsubscribed", handleFilterUnsubscribed)
	http.HandleFunc("/queryUnsubscribed", handleQueryUnsubscribed)
	http.ListenAndServe(":9090", nil)
}
