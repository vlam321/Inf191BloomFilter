package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/vlam321/Inf191BloomFilter/bloomDataGenerator"
	"github.com/vlam321/Inf191BloomFilter/databaseAccessObj"
	"github.com/vlam321/Inf191BloomFilter/payload"

	"github.com/cyberdelia/go-metrics-graphite"
	metrics "github.com/rcrowley/go-metrics"
)

var endpoint string
var host string

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

func makeMap(emails []string) map[string]bool {
	quickMap := make(map[string]bool)
	for i := range emails {
		quickMap[emails[i]] = true
	}
	return quickMap
}

// checkResult takes in the expected and actual values and
// calculate the hit and miss ratio and sends the data to
// graphite
func checkResult(unsubbed, subbed map[int][]string, res []string) {
	unsubbedMap := makeMap(unsubbed[0])
	subbedMap := makeMap(subbed[0])
	hit := 0
	miss := 0
	for i := range res {
		if ok := (unsubbedMap[res[i]] && !subbedMap[res[i]]); ok {
			hit += 1
		} else {
			miss += 1
		}
	}
	if len(unsubbedMap) > len(res) {
		miss += len(unsubbedMap) - len(res)
	}
	metrics.GetOrRegisterGauge("result.hit", nil).Update(int64(hit))
	metrics.GetOrRegisterGauge("result.miss", nil).Update(int64(miss))
}

// attackBloomFilter hit endpoint with test data
func attackBloomFilter(dao *databaseAccessObj.Conn, expectedTrues, expectedFalse int, endpoint string, userID int) {

	unsubbed := dao.SelectRandSubset(userID, expectedTrues)
	subbed := bloomDataGenerator.GenData(1, expectedFalse, expectedFalse+1)

	var dataSum []string
	dataSum = append(dataSum, unsubbed[userID]...)
	dataSum = append(dataSum, subbed[0]...)

	pyld := Payload{userID, dataSum}
	jsn := conv2Json(pyld)

	// membershipEndpoint := "http://" + endpoint + ":9090/filterUnsubscribed"
	// log.Println(membershipEndpoint)

	start := time.Now()
	client := &http.Client{}
	//_, err := http.Post(endpoint, "application/json; charset=utf-8", bytes.NewBuffer(jsn))
	r, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsn))
	if err != nil {
		log.Printf("Error in post request: %v\n", err)
		return
	}
	r.Header.Set("userid", strconv.Itoa(pyld.UserId))
	r.Header.Add("Content-Type", "application/json; charset=utf-8")
	res, err := client.Do(r)
	latency := time.Since(start).Nanoseconds() / 1000000

	if err != nil {
		log.Printf("Error in post request: %v\n", err)
		return
	}

	log.Printf("Sent request to filter with payload size of %d emails (expected reponse size = %d emails).", len(dataSum), expectedTrues)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		return
	}

	//var result map[int][]string
	var result payload.Payload
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error unmarshaling body: %v\n", err)
		return
	}

	metrics.GetOrRegisterGauge(fmt.Sprintf("%s.request.latency", host), nil).Update(int64(latency))
	metrics.GetOrRegisterGauge(fmt.Sprintf("%s.request.trues", host), nil).Update(int64(len(unsubbed[userID])))
	metrics.GetOrRegisterGauge(fmt.Sprintf("%s.request.falses", host), nil).Update(int64(len(dataSum) - len(unsubbed[userID])))
	// checkResult(unsubbed, subbed, result.Emails)
}

// sendRequest attackBloomFilter every ms
func sendRequest(dao *databaseAccessObj.Conn, ms int32, endpoint string, userID int) {
	ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
	for _ = range ticker.C {
		attackBloomFilter(dao, 500, 500, endpoint, userID)
	}
}

func main() {
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:2003", os.Getenv("GRAPHITE_IP")))
	host, _ = os.Hostname()

	go graphite.Graphite(metrics.DefaultRegistry, 10e9, "metrics", addr)

	endpoint = os.Getenv("ENDPOINT")
	userID, err := strconv.Atoi(os.Getenv("USERID"))
	if err != nil {
		log.Printf("Client simulator: %v\n", err.Error())
	}

	log.Printf("ATTACKING ROUTER @ %s\n", endpoint)
	go sendRequest(dao, 1, endpoint, userID)

	http.ListenAndServe(":9091", nil)
}
