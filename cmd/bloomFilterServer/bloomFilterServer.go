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
	var buff []byte
	var payload Payload

	err := json.NewDecoder(r.Body).Decode(&buff)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(buff, &payload)
	checkErr(err)

	fmt.Println(payload)

	bf := bloomManager.New()
	emails := bf.GetArrayOfUnsubscribedEmails(payload.Emails)
	filteredEmails := Result{emails}
	// buff2 := new(bytes.Buffer)
	// data, err := json.Marshal(filteredEmails)
	// checkErr(err)

	js, err := json.Marshal(filteredEmails)
	checkErr(err)

	fmt.Println(js)
	w.Write(js)
	//return the result in whatever format use http.response
	//look into resturing body
	//encode to type.buffer
}

func main() {
	http.HandleFunc("/filterUnsubscribed", handleFilterUnsubscribed)
	http.ListenAndServe(":9090", nil)
}
