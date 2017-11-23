/* This file will contain functionalities to host the main service
in which the client will make requests to. Use net/http to create
API endpoints to access the following functionalities:
	- UpdateBloomFilter
	- RepopulateBloomFilter
	- GetArrayOfUnsubscribedEmails
	- etc. (Steph
		add more functionalities
		here and implement this
		after you've completed your
		current tasks)
*/
package main

import (
	"Inf191BloomFilter/bloomManager"
	"encoding/json"
	"fmt"
	"net/http"
)

//struct that takes an int and a list of emails
type Payload struct {
	UserId int
	Emails []string
}

type Result struct {
	Trues []string
}

//handleUpdate will update the respective bloomFilter
func handleUpdate(r http.ResponseWriter, req *http.Request) {
	//	Use UpdateBloomFilter() here
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func handleFilterUnsubscribed(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request: %v %v %v\n", r.Method, r.URL, r.Proto)
	var buff []byte
	var payload Payload

	err := json.NewDecoder(r.Body).Decode(&buff)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(buff, &payload)
	checkErr(err)

	bf := bloomManager.New()
	emails := bf.GetArrayOfUnsubscribedEmails(payload.Emails)
	filteredEmails := Result{emails}

	js, err := json.Marshal(filteredEmails)
	checkErr(err)

	w.Write(js)
}

func main() {
	http.HandleFunc("/filterUnsubscribed", handleFilterUnsubscribed)
	http.ListenAndServe(":9090", nil)
}
