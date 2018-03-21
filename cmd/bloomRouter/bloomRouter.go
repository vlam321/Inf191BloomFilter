package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"strconv"

	"github.com/fatih/structs"
	"github.com/spf13/viper"
)

// BloomServerIPs struct holding ips of each bloom filter server; retrieved through getBloomFilterIPs()
type BloomServerIPs struct {
	BloomFilterServer1  string
	BloomFilterServer2  string
	BloomFilterServer3  string
	BloomFilterServer4  string
	BloomFilterServer5  string
	BloomFilterServer6  string
	BloomFilterServer7  string
	BloomFilterServer8  string
	BloomFilterServer9  string
	BloomFilterServer10 string
}

type BloomContainerNames struct {
	BloomFilterContainer1  string
	BloomFilterContainer2  string
	BloomFilterContainer3  string
	BloomFilterContainer4  string
	BloomFilterContainer5  string
	BloomFilterContainer6  string
	BloomFilterContainer7  string
	BloomFilterContainer8  string
	BloomFilterContainer9  string
	BloomFilterContainer10 string
}

var bloomServerIPs BloomServerIPs
var bloomContainerNames BloomContainerNames
var routes map[int]string
var re *regexp.Regexp

func retrieveEndpoint(userid int) string {
	var endpoint string
	if viper.GetString("host") == "ecs" {
		if os.Getenv("SKIP_FILTER") == "true" {
			endpoint = "http://" + os.Getenv(routes[userid]) + ":9090/queryUnsubscribed"
		} else {
			endpoint = "http://" + os.Getenv(routes[userid]) + ":9090/filterUnsubscribed"
		}
	} else {
		endpoint = "http://" + viper.GetString("dockerIP") + ":" + routes[userid] + "/filterUnsubscribed"
	}
	return endpoint
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
	// read request data
	/*
		bbytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error: Unable to read request data. %v\n", err)
			return
		}
	*/
	userid, err := strconv.Atoi(r.Header.Get("userid"))
	if err != nil {
		log.Printf("String conversion error. %v\n", err.Error())
	}

	/*
		vals := strings.Split(string(bbytes), ",")
		var userid int
		for i := range vals {
			if strings.Contains(vals[i], "UserId") {
				userid, err = strconv.Atoi(re.FindString(vals[i]))
				if err != nil {
					log.Printf("strconv error: %v\n", err)
					return
				}
				break
			}
		}
	*/
	/*
		// unmarshal payload
		var pl payload.Payload
		err = json.Unmarshal(bbytes, &pl)
		if err != nil {
			log.Printf("Error: Unable to unmarshal Payload. %v\n", err)
			return
		}

		// determine endpoint based on host
		endpoint := retrieveEndpoint(pl.UserId)
	*/

	endpoint := retrieveEndpoint(userid)
	log.Printf("Request sent to: %s\n", endpoint)
	http.Redirect(w, r, endpoint, http.StatusTemporaryRedirect)

	/*
		// make request to endpoint
		res, _ := http.Post(endpoint, "application/json; charset=utf-8", bytes.NewBuffer(bbytes))
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Router: error reading response from bloom filter. %v\n", err)
		}
	*/
	// w.Write(body)

}

// getMyIP() retrieve IP on local host
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

// getBloomFilterIPs() retrieve IPs of each bloom filter server and store in bloomServerIPs
func getBloomFilterIPs() error {
	viper.SetConfigName("bfIPConf")
	viper.AddConfigPath("settings")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	if viper.GetString("host") == "ecs" {
		log.Printf("host: ecs")
		err = viper.Unmarshal(&bloomContainerNames)
		if err != nil {
			return err
		}
	} else {
		log.Printf("host:docker")
		err = viper.Unmarshal(&bloomServerIPs)
		if err != nil {
			return err
		}
	}
	return nil
}

func mapRouter(bloomFilterIPs BloomServerIPs) {
	re = regexp.MustCompile("[0-9]+")
	routes = make(map[int]string)
	if viper.GetString("host") == "ecs" {
		containerNames := structs.Values(bloomContainerNames)
		for i := range containerNames {
			routes[i] = containerNames[i].(string)
		}
	} else {
		bloomIPs := structs.Values(bloomFilterIPs)
		for i := range bloomIPs {
			routes[i] = bloomIPs[i].(string)
		}
	}
	for k, v := range routes {
		log.Printf("key: %v	 value: %v\n", k, v)
	}
}

func main() {
	log.Printf("RETRIEVING BLOOM FILTER IP'S")
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
