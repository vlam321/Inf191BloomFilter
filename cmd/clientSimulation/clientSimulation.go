package main

import (
	"Inf191BloomFilter/bloomDataGenerator"
	"Inf191BloomFilter/databaseAccessObj"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cyberdelia/go-metrics-graphite"
	metrics "github.com/rcrowley/go-metrics"
)

const membershipEndpoint = "http://localhost:9090/filterUnsubscribed"

type Payload struct {
	UserId int
	Emails []string
}

// conv2Json converts payload input into JSON
func conv2Json(payload Payload) []byte {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error json marshaling: %v\n", err)
		return nil
	}
	return data
}

// attackBloomFilter hit endpoint with test data 
func attackBloomFilter(dao *databaseAccessObj.Conn) {
	unsubbed := dao.SelectRandSubset(0, 1000)
	subbed := bloomDataGenerator.GenData(1, 100, 200)

	var dataSum []string
	dataSum = append(dataSum, unsubbed[0]...)
	dataSum = append(dataSum, subbed[0]...)

	pyld := Payload{0, dataSum}
	jsn := conv2Json(pyld)

	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", bytes.NewBuffer(jsn))

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		return
	}

	var result map[int][]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error unmarshaling body: %v\n", err)
		return
	}
	log.Printf("Success: %d emails returned\n", len(result[0]))
}

// sendRequest attackBloomFilter every ms
func sendRequest(dao *databaseAccessObj.Conn, ms int32) {
	ticker := time.NewTicker(time.Duration(ms)*time.Millisecond)
	for _ = range ticker.C {
		attackBloomFilter(dao)
	}
}

func main() {
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	addr, _ := net.ResolveTCPAddr("tcp", "192.168.99.100:2003")
	go graphite.Graphite(metrics.DefaultRegistry, 10e9, "metrics", addr)
	go sendRequest(dao, 1000)
	http.ListenAndServe(":9091", nil)
}
