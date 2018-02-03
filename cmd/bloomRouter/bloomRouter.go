package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fatih/structs"
	"github.com/spf13/viper"
)

type Payload struct {
	UserId int
	Emails []string
}

type BloomServerIPs struct {
	BloomFilterServer1 string
	BloomFilterServer2 string
	BloomFilterServer3 string
	BloomFilterServer4 string
	BloomFilterServer5 string
}

var bloomServerIPs BloomServerIPs
var routes map[int]string

func handleRoute(w http.ResponseWriter, r *http.Request) {
	bbytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error: Unable to read request data. %v\n", err.Error())
		return
	}

	var pyld Payload
	err = json.Unmarshal(bbytes, &pyld)
	if err != nil {
		log.Printf("Error: Unable to unmarshal Payload. %v\n", err.Error())
		return
	}
	var endpoint string
	if viper.GetString("host") == "ec2" {
		endpoint = "http://" + routes[pyld.UserId] + ":9090/filterUnsubscribed"
	} else {
		endpoint = "http://" + viper.GetString("dockerIP") + ":" + routes[pyld.UserId] + "/filterUnsubscribed"
	}
	log.Printf("Request sent to: %s\n", endpoint)

	res, _ := http.Post(endpoint, "application/json; charset=utf-8", bytes.NewBuffer(bbytes))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Router: error reading response from bloom filter. %v\n", err.Error())
	}
	w.Write(body)
}

func getMyIP() (myIP string, err error) {
	resp, err := http.Get("http://checkip.amazonaws.com/")
	if err != nil {
		return "x.x.x.x", errors.New("Unable to find IP.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "x.x.x.x", errors.New("Unable to find IP.")
	}
	return string(body[:]), nil
}

func getBloomFilterIPs() error {
	viper.SetConfigName("bfIPConf")
	viper.AddConfigPath("settings")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&bloomServerIPs)
	if err != nil {
		return err
	}

	return nil
}

func mapRouter(bloomFilterIPs BloomServerIPs) {
	routes = make(map[int]string)
	bloomIPs := structs.Values(bloomFilterIPs)
	for i := range bloomIPs {
		routes[i] = bloomIPs[i].(string)
	}
}

func main() {
	err := getBloomFilterIPs()
	if err != nil {
		log.Println(err)
	}
	log.Printf("SUCCESSFULLY PARSED BLOOM SERVER IPS.")

	mapRouter(bloomServerIPs)
	log.Printf("SUCCESSFULLY MAPPED BLOOM SERVER IPS.")

	http.HandleFunc("/filterUnsubscribed", handleRoute)
	http.ListenAndServe(":9090", nil)
}
