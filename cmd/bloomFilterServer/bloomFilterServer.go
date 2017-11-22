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
	"strconv"
)

//struct that takes an int and a list of emails
type Payload struct {
	UserId int
	Emails []string
}

//handleUpdate will update the respective bloomFilter
func handleUpdate(r http.ResponseWriter, req *http.Request) {
	//	Use UpdateBloomFilter() here
}

func handleFilterUnsubscribed(w http.ResponseWriter, r *http.Request) {
	var p Payload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var idEmail string
	var arrayOfEmails []string
	for i := range p.Emails {
		idEmail = (strconv.Itoa(p.UserId) + p.Emails[i])
		arrayOfEmails = append(arrayOfEmails, idEmail)
	}

	bf := bloomManager.New()
	var c []string
	c = bf.GetArrayOfUnsubscribedEmails(arrayOfEmails)
	fmt.Fprint
	//return the result in whatever format use http.response
	//look into resturing body
	//encode to type.buffer
}

func main() {
	http.HandleFunc("/filterUnsubscribed", handleFilterUnsubscribed)
	http.ListenAndServe(":9090", nil)
}
