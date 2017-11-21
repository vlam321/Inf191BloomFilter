/* This file will contain functionality for a simulated
client. The client is responsible for makine http requests
to the dbServer and bloomFilterServer
*/

package bloomFilterServer

import (
	//"encoding/json"
	//"bytes"
	"fmt"
	//"net/http"
	"testing"
	"Inf191BloomFilter/bloomDataGenerator"
	"Inf191BloomFilter/databaseAccessObj"
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
	// var payload Payload
	var dataSum []string
	buff := new(bytes.Buffer)

	// Generate random id_email pairs (positives) and save it in a var
	inDB := bloomDataGenerator.GenData(1, 100, 200)

	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	// Insert new data in the db
	dao.Insert(inDB)

	// Call BF server to update the bit array
	// res, err := http.Get(update_bit_array_endpoint)
	// checkErr(err)

	// Generate more raandom id_email pairs (negatives) and save it ina var
	notInDB := bloomDataGenerator.GenData(1, 50, 100)

	for userid, emails := range inDB {
		dataSum = append(emails, notInDB[userid]...)
	}

	// Put values into payload to be sent to the server later
	payload = Payload{0, dataSum}

	// convert to json
	data, err := json.Marshal(payload)
	checkErr(err)

	// encode to bytes buffer
	err = json.NewEncoder(buff).Encode(data)
	checkErr(err)

	// requst for true memberships
	res, _ := http.Post(membership_endpoint, "application/json; charset=utf-8", buff)
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
