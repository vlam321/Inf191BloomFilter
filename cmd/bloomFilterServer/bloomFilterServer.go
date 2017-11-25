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
type Payload struct {
	UserId int
	Emails []string
}

//Global variable 
//The bloom filter for this server
var bf = bloomManager.New()


//handleUpdate will update the respective bloomFilter
func handleUpdate(r http.ResponseWriter, req *http.Request) {
	bf.UpdateBloomFilter()

}

//checkErr checks fro errors in the decoded text and encoded text... 
func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}


func handleFilterUnsubscribed(w http.ResponseWriter, r *http.Request) {
	var buff []byte
	var payload Payload
	//Result struct made to carry the result of unsuscribed emails 
	type Result struct {
	Trues []string
	}

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

	js, err := json.Marshal(filteredEmails)
	checkErr(err)

	fmt.Println(js)
	w.Write(js)
}

func main() {
	http.HandleFunc("/filterUnsubscribed", handleFilterUnsubscribed)
	http.HandleFunc("/update", handleUpdate)
	
	http.ListenAndServe(":9090", nil)
}
