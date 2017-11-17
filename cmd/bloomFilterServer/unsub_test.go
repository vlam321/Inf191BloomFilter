/* This file will contain functionality for a simulated
client. The client is responsible for makine http requests
to the dbServer and bloomFilterServer
*/

package main

import (
	"encoding/json"
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"Inf191BloomFilter/bloomDataGenerator"
)

const membership_endpoint = "http://localhost:9090/members"

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Payload struct {
	UserId int
	Emails []string
}

func TestUnsub(t *testing.T) {
	/*
		1. gen random user_email pairs, and save to a var
		2. insert the user_email pairs intot the db
		3. call BF server to update the bit array
		4. gen more data and concatenate the saved pairs to the new ones
		5. run the concatenate user_email pairs against the BF server, and 
		   make sure that only the first saved ones are in the response
	*/
	payload := Payload{1, []string{"sodfd", "fdsafasd"}}
	buff := new(bytes.Buffer)

	data, err := json.Marshal(payload)
	checkErr(err)

	err = json.NewEncoder(buff).Encode(data)
	checkErr(err)

	res, _ := http.Post(membership_endpoint, "application/json; charset=utf-8", buff)
	fmt.Println(res)
}

func TestNewDataAdded(t *testing.T) {
	/*
		1. gen rand data, store it in true dataset var and insert into db
		2. run data against BF and make sure returns empty slice
		3. insert data into db using doa
		4. run request for updating BF bit array
		5. rerun data again BF and make sure len(res) == len(data)
	*/
}


func TestDeleteData(t * testing.T) {
	/*
		1. use dao to grab some data pairs, store in a var
		2. use dao to remove these pairs from the db
		3. call BF server to update the bit array
		4. Run the saved pairs against the BF server and make sure the get
		an empty array
	*/
}
