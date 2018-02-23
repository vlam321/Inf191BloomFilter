package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/vlam321/Inf191BloomFilter/bloomDataGenerator"
	"github.com/vlam321/Inf191BloomFilter/databaseAccessObj"

	"github.com/cyberdelia/go-metrics-graphite"
	metrics "github.com/rcrowley/go-metrics"
)

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
func checkResult(unsubbed, subbed, res map[int][]string) {
	unsubbedMap := makeMap(unsubbed[0])
	subbedMap := makeMap(subbed[0])
	hit := 0
	miss := 0
	for i := range res[0] {
		if ok := (unsubbedMap[res[0][i]] && !subbedMap[res[0][i]]); ok {
			hit += 1
		} else {
			miss += 1
		}
	}
	if len(unsubbedMap) > len(res[0]) {
		miss += len(unsubbedMap) - len(res[0])
	}
	metrics.GetOrRegisterGauge("result.hit", nil).Update(int64(hit))
	metrics.GetOrRegisterGauge("result.miss", nil).Update(int64(miss))
}

// attackBloomFilter hit endpoint with test data
func attackBloomFilter(dao *databaseAccessObj.Conn, expectedTrues, expectedFalse int, routerIp string) {
	unsubbed := dao.SelectRandSubset(0, expectedTrues)
	subbed := bloomDataGenerator.GenData(1, expectedFalse, expectedFalse+501)
	var dataSum []string
	dataSum = append(dataSum, unsubbed[0]...)
	dataSum = append(dataSum, subbed[0]...)

	pyld := Payload{0, dataSum}
	jsn := conv2Json(pyld)

	membershipEndpoint := "http://" + routerIp + ":9090/filterUnsubscribed"
	log.Println(membershipEndpoint)
	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", bytes.NewBuffer(jsn))

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		return
	}

	log.Printf("Sent request to filter with payload size of %d emails (expected reponse size = %d emails).", len(dataSum), expectedTrues)

	var result map[int][]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error unmarshaling body: %v\n", err)
		return
	}
	metrics.GetOrRegisterGauge("request.hit", nil).Update(int64(len(result[0])))
	metrics.GetOrRegisterGauge("request.miss", nil).Update(int64(len(dataSum) - len(result[0])))

	checkResult(unsubbed, subbed, result)
}

// sendRequest attackBloomFilter every ms
func sendRequest(dao *databaseAccessObj.Conn, ms int32, routerIp string) {
	ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
	for _ = range ticker.C {
		attackBloomFilter(dao, 2000, 500, routerIp)
	}
}

func main() {
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	addr, _ := net.ResolveTCPAddr("tcp", "192.168.99.100:2003")
	go graphite.Graphite(metrics.DefaultRegistry, 10e9, "metrics", addr)
	routerIp := os.Getenv("ROUTER_IP")
	log.Printf("ATTACKING ROUTER @ %s", routerIp)
	go sendRequest(dao, 2000, routerIp)
	http.ListenAndServe(":9091", nil)
}
