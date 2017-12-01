/* This file will contain functionalities to host the main service
in which the client will make requests to. Use net/http to create
API endpoints to access the following functionalities:
	- UpdateBloomFilter
		- AddToBloomFilter
		- RepopulateBloomFilter
	- GetArrayOfUnsubscribedEmails
*/
package main

import (
	"Inf191BloomFilter/bloomManager"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cyberdelia/go-metrics-graphite"
	metrics "github.com/rcrowley/go-metrics"
)

type payload struct {
	UserId int
	Emails []string
}

//The bloom filter for this server
var bf *bloomManager.BloomFilter

//handleUpdate will update the respective bloomFilter
//server will keep track of when the last updated time is. to call
//update every _ time.
func handleUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request: %v %v %v\n", r.Method, r.URL, r.Proto)
	bf.RepopulateBloomFilter()
}

//checkErr checks fro errors in the decoded text and encoded text...
func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func handleMetric(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("1")
		return
	}
	u := r.Form["user"]
	if len(u) == 0 {
		log.Printf("2")
		return
	}

	uid, err := strconv.Atoi(u[0])
	if err != nil {
		log.Printf("3")
		return
	}
	metrics.GetOrRegisterGauge("userid.gauge", nil).Update(int64(uid))
	metrics.GetOrRegisterCounter("userid.counter", nil).Inc(1)
	if err != nil {
		log.Printf("5")
		return
	}

	log.Printf("user id  = %d\n", uid)
}

//handleFilterUnsubscripe handles when the client requests to recieve
//those emails that are unsubscribed therefore, IN the database.
func handleFilterUnsubscribed(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request: %v %v %v\n", r.Method, r.URL, r.Proto)
	var buff []byte
	var payload payload

	//Result struct made to carry the result of unsuscribed emails
	type Result struct {
		Trues []string
	}

	// log.Printf("%v", r.Body)
	//decode byets from request body
	err := json.NewDecoder(r.Body).Decode(&buff)

	//check for error; if error write 404
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
		// http.Error(w, err.Error(), 400)
		// return
	}

	// converts the decoded result to a payload struct
	err = json.Unmarshal(buff, &payload)
	if err != nil {
		log.Printf("error unmashaling payload %v %s\n", err, string(buff))
		return
	}
	var emailInputs []string
	for i := range payload.Emails {
		emailInputs = append(emailInputs, strconv.Itoa(payload.UserId)+"_"+payload.Emails[i])
	}

	//uses bloomManager to get the result of unsubscribed emails
	//puts them in struct, result
	emails := bf.GetArrayOfUnsubscribedEmails(emailInputs)
	filteredEmails := Result{emails}

	//convert result to json
	js, err := json.Marshal(filteredEmails)
	checkErr(err)

	//write back to client
	w.Write(js)
}

func timeManager() {
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for t := range ticker.C {
			//Call update bloom filter
			bf.UpdateBloomFilter()
			fmt.Println("Bloom Filter Updated at: ", t.Format("2006-01-02 3:4:5 PM"))

		}
	}()
	//Figure out how to run without sleep?
	//use go forever? an option
	time.Sleep(time.Second * 10)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func setBloomFilter(bitArraySize, numHashFunc uint) {
	bf = bloomManager.New(bitArraySize, numHashFunc)
}

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "192.168.99.100:2003")
	go graphite.Graphite(metrics.DefaultRegistry, 10e9, "metrics", addr)
	bitArraySize, _ := strconv.ParseUint(os.Args[1], 10, 64)
	numHashFunc, _ := strconv.ParseUint(os.Args[2], 10, 64)
	setBloomFilter(uint(bitArraySize), uint(numHashFunc))

	timeManager()
	http.HandleFunc("/filterUnsubscribed", handleFilterUnsubscribed)
	http.HandleFunc("/update", handleUpdate)
	http.HandleFunc("/metric", handleMetric)

	http.ListenAndServe(":9090", nil)
}
