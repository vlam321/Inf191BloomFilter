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
	"encoding/json"
	"fmt"
	"net/http"
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

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func handleFilterUnsubscribed(w http.ResponseWriter, r *http.Request) {
	var buff []byte
	var payload interface{}

	err := json.NewDecoder(r.Body).Decode(&buff)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	payload = json.Unmarshal(buff, payload)

	fmt.Println(payload)
	/*
		var idEmail string
		var arrayOfEmails []string
		for i := range p.Emails {
			idEmail = (strconv.Itoa(p.UserId) + p.Emails[i])
			arrayOfEmails = append(arrayOfEmails, idEmail)
		}

		bf := bloomManager.New()
		var c []string
		c = bf.GetArrayOfUnsubscribedEmails(arrayOfEmails)

		buff := new(bytes.Buffer)
		data, err := json.Marshal(c)
		checkErr(err)
		fmt.Println("1")
		err = json.NewEncoder(buff).Encode(&data)
		checkErr(err)
		fmt.Println("2")
		fmt.Println(data)

		w.Write(data)
	*/
	//return the result in whatever format use http.response
	//look into resturing body
	//encode to type.buffer
}

func main() {
	http.HandleFunc("/filterUnsubscribed", handleFilterUnsubscribed)
	http.ListenAndServe(":9090", nil)
}
