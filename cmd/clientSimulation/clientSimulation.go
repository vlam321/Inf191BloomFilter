package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/vlam321/Inf191BloomFilter/bloomDataGenerator"
	"github.com/vlam321/Inf191BloomFilter/databaseAccessObj"

	"github.com/cyberdelia/go-metrics-graphite"
	metrics "github.com/rcrowley/go-metrics"
)

var endpoint string
var host string
var trueValues []string
var falseValues []string

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

/* LEGACY CODE
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
*/

// Takes a string array and shuffle the value
func Shuffle(a []string) {
	rand.Seed(time.Now().UnixNano())
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// takes a string array, shuffles it's value and
// returns a random subset with a specific amount of values
func getRandSubset(a []string, n int) []string {
	Shuffle(a)
	return a[0:n]

}

// attackBloomFilter hit endpoint with test data
func attackBloomFilter(expectedTrues, expectedFalse int, endpoint string, userID int) {

	// Get random subset of trueValues and falseValues
	unsubbed := getRandSubset(trueValues, expectedTrues)
	subbed := getRandSubset(falseValues, expectedFalse)

	var dataSum []string
	dataSum = append(dataSum, unsubbed...)
	dataSum = append(dataSum, subbed...)

	pyld := Payload{userID, dataSum}
	jsn := conv2Json(pyld)

	client := &http.Client{}

	r, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsn))
	if err != nil {
		log.Printf("Error in post request: %v\n", err)
		return
	}
	r.Header.Set("userid", strconv.Itoa(pyld.UserId))
	r.Header.Add("Content-Type", "application/json; charset=utf-8")
	start := time.Now()
	_, err = client.Do(r)
	latency := time.Since(start).Nanoseconds() / 1000000

	if err != nil {
		log.Printf("Error in post request: %v\n", err)
		return
	}

	log.Printf("Sent request to filter with payload size of %d emails (expected reponse size = %d emails).", len(dataSum), expectedTrues)
	metrics.GetOrRegisterGauge(fmt.Sprintf("%s.request.latency", host), nil).Update(int64(latency))
	metrics.GetOrRegisterGauge(fmt.Sprintf("%s.request.trues", host), nil).Update(int64(len(unsubbed)))
	metrics.GetOrRegisterGauge(fmt.Sprintf("%s.request.falses", host), nil).Update(int64(len(subbed)))

	/* LEGACY CODE
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

	checkResult(unsubbed, subbed, result.Emails)
	*/
}

// sendRequest attackBloomFilter every ms
func sendRequest(ms int32, endpoint string, userID int) {
	ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
	for _ = range ticker.C {
		attackBloomFilter(500, 500, endpoint, userID)
	}
}

func main() {

	// Get neccessary environment variables
	endpoint = os.Getenv("ENDPOINT")
	userID, err := strconv.Atoi(os.Getenv("USERID"))
	if err != nil {
		log.Printf("Client simulator: %v\n", err.Error())
	}

	// Open new DB connection
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()

	// Grab 5000 known unsubbed emails (trueValues) from DB
	// Generate 5000 new random emails to represent subbed emails (falseValues)
	trueValues = dao.SelectRandSubset(userID, 5000)[userID]
	falseValues = bloomDataGenerator.GenData(1, 5000, 5001)[0]

	// Assign dependent values to run graphite
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:2003", os.Getenv("GRAPHITE_IP")))
	host, _ = os.Hostname()

	// Run graphite
	go graphite.Graphite(metrics.DefaultRegistry, 10e9, "metrics", addr)

	// Send Request to bloom router
	log.Printf("ATTACKING ROUTER @ %s\n", endpoint)
	go sendRequest(1, endpoint, userID)

	http.ListenAndServe(":9091", nil)
}
