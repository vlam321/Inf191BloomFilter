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
	"Inf191BloomFilter/bloomManager"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
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
	checkErr(err)
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

//updateBloomFilterBackground manages the periodic updates of the bloom
//filter. Update calls repopulate, creating a new updated bloom filter
func updateBloomFilterBackground() {
	//Set new ticker to every 2 seconds
	ticker := time.NewTicker(time.Second * 2)

	for t := range ticker.C {
		//Call update bloom filter
		bf.UpdateBloomFilter()
		fmt.Println("Bloom Filter Updated at: ", t.Format("2006-01-02 3:4:5 PM"))
	}
}

func setBloomFilter(bitArraySize, numHashFunc uint) {
	bf = bloomManager.New(bitArraySize, numHashFunc)
}

func main() {

	bitArraySize, _ := strconv.ParseUint(os.Args[1], 10, 64)
	numHashFunc, _ := strconv.ParseUint(os.Args[2], 10, 64)
	setBloomFilter(uint(bitArraySize), uint(numHashFunc))

	//Run go routine to make periodic updates
	//Runs until the server is stopped
	go updateBloomFilterBackground()
	http.HandleFunc("/filterUnsubscribed", handleFilterUnsubscribed)
	http.HandleFunc("/update", handleUpdate)

	http.ListenAndServe(":9090", nil)
}
