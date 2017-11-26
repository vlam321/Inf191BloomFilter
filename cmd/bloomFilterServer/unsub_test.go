/* This file will contain functionality for a simulated
client. The client is responsible for makine http requests
to the dbServer and bloomFilterServer
*/

package bloomFilterServer

import (
	"Inf191BloomFilter/bloomDataGenerator"
	"Inf191BloomFilter/databaseAccessObj"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dsn = "bloom:test@/unsubscribed"
const membershipEndpoint = "http://localhost:9090/filterUnsubscribed"
const updateEndpoint = "http://localhost:9090/update"

type Payload struct {
	UserId int
	Emails []string
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Result struct {
	Trues []string
}

func TestUnsub(t *testing.T) {
	dao := databaseAccessObj.New(dsn)
	// Clear out values in table
	dao.Clear()

	var dataSum []string
	buff := new(bytes.Buffer)
	var payload Payload

	// Generate random id_email pairs (positives) and save it in a var
	// Increasing these values may produce false positives
	inDB := bloomDataGenerator.GenData(1, 10, 20)

	// Insert new data in the db
	dao.Insert(inDB)

	// Call BF server to update the bit array
	_, err := http.Get(updateEndpoint)
	checkErr(err)

	// Generate more raandom id_email pairs (negatives) and save it ina var
	notInDB := bloomDataGenerator.GenData(1, 50, 100)

	// Concatenate the two slices
	for userid, emails := range inDB {
		dataSum = append(emails, notInDB[userid]...)
	}

	fmt.Println(len(dataSum))
	// Put values into payload to be sent to the server later
	payload = Payload{0, dataSum}

	// Convert to json
	data, err := json.Marshal(payload)
	checkErr(err)

	// Encode to bytes buffer
	err = json.NewEncoder(buff).Encode(data)
	checkErr(err)

	// Request for members that exist in DB
	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", buff)

	var arr Result

	// Read the result
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	// converts the decoded result back to a Result struct
	err = json.Unmarshal(body, &arr)
	checkErr(err)

	fmt.Printf("%d == %d\n", len(arr.Trues), len(inDB[0]))
	// checks that only values in DB are returned
	assert.True(t, len(arr.Trues) == len(inDB[0]))
}

func TestNewUnsubscribes(t *testing.T) {
	dao := databaseAccessObj.New(dsn)
	// Clear values in tables for clean test
	dao.Clear()

	buff := new(bytes.Buffer)

	// Increasing these values may produce false positives
	dataSet := bloomDataGenerator.GenData(1, 10, 20)
	payload := Payload{0, dataSet[0]}

	data, err := json.Marshal(payload)
	checkErr(err)

	err = json.NewEncoder(buff).Encode(data)
	checkErr(err)

	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", buff)

	var arr Result

	// Read the result
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	// converts the decoded result back to a Result struct
	err = json.Unmarshal(body, &arr)
	checkErr(err)

	fmt.Printf("%d != %d\n", len(arr.Trues), len(dataSet[0]))
	assert.True(t, len(arr.Trues) == 0)

	// Insert the true values into the database
	dao.Insert(dataSet)

	// Update the bloom filter bit array
	_, err = http.Get(updateEndpoint)
	checkErr(err)
	err = json.NewEncoder(buff).Encode(data)
	checkErr(err)

	res, _ = http.Post(membershipEndpoint, "application/json; charset=utf-8", buff)

	body, err = ioutil.ReadAll(res.Body)
	checkErr(err)

	var arr2 Result
	// converts the decoded result back to a Result struct
	err = json.Unmarshal(body, &arr2)
	checkErr(err)

	// Check again to see if bit array is updated
	fmt.Printf("%d == %d\n", len(arr2.Trues), len(dataSet[0]))
	assert.True(t, len(arr2.Trues) == len(dataSet[0]))
}

func TestResubscribed(t *testing.T) {
	dao := databaseAccessObj.New(dsn)

	// Clear tables for clean tests
	dao.Clear()

	/*
		1. use dao to grab some data pairs, store in a var
		2. use dao to remove these pairs from the db
		3. call BF server to update the bit array
		4. Run the saved pairs against the BF server and make sure the get an empty array
	*/

	var allData []string
	dataSet := bloomDataGenerator.GenData(1, 10, 20)
	extra := bloomDataGenerator.GenData(1, 10, 20)

	allData = append(allData, dataSet[0]...)
	allData = append(allData, extra[0]...)

	dataSum := make(map[int][]string)
	dataSum[0] = allData

	dao.Insert(dataSum)

	_, err := http.Get(updateEndpoint)
	checkErr(err)

	// Remove the extra from the database
	dao.Delete(extra)

	// update the bloom filter bit array again
	_, err = http.Get(updateEndpoint)
	checkErr(err)

	buff := new(bytes.Buffer)
	payload := Payload{0, dataSum[0]}

	data, err := json.Marshal(payload)
	checkErr(err)

	err = json.NewEncoder(buff).Encode(data)
	checkErr(err)

	// get the memberships using yall the data
	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", buff)

	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	var arr Result
	// converts the decoded result back to a Result struct
	err = json.Unmarshal(body, &arr)
	checkErr(err)

	fmt.Printf("%d != %d", len(dataSum[0]), len(arr.Trues))
	assert.False(t, len(dataSum[0]) == len(arr.Trues))

	// Grab some of these data from the database

}
