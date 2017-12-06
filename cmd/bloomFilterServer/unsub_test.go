/* This file will contain functionality for a simulated
client. The client is responsible for makine http requests
to the dbServer and bloomFilterServer
*/

package main

import (
	"Inf191BloomFilter/bloomDataGenerator"
	"Inf191BloomFilter/databaseAccessObj"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const membershipEndpoint = "http://localhost:9090/filterUnsubscribed"
const updateEndpoint = "http://localhost:9090/update"

/*
type Payload struct {
	UserId int
	Emails []string
}
*/

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Result struct {
	Trues []string
}

func TestUnsub(t *testing.T) {
	dao := databaseAccessObj.New()
	// Clear out values in table
	dao.Clear()

	var dataSum []string
	var pyld Payload

	// Generate random id_email pairs (positives) and save it in a var
	// Increasing these values may produce false positives
	inDB := bloomDataGenerator.GenData(1, 1000, 2000)

	// Insert new data in the db
	dao.Insert(inDB)

	// Call BF server to update the bit array
	_, err := http.Get(updateEndpoint)
	checkErr(err)

	// Generate more raandom id_email pairs (negatives) and save it ina var
	notInDB := bloomDataGenerator.GenData(1, 1000, 2000)

	// Concatenate the two slices
	for userid, emails := range inDB {
		dataSum = append(emails, notInDB[userid]...)
	}

	fmt.Println("Total ID:Email pairs inserted: ", len(dataSum))
	// Put values into Payload to be sent to the server later
	pyld = Payload{0, dataSum}

	// Convert to json
	data, err := json.Marshal(pyld)
	checkErr(err)

	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", bytes.NewBuffer(data))

	// Request for members that exist in D
	// Read the result
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var result map[int][]string

	// converts the decoded result back to a Result struct
	err = json.Unmarshal(body, &result)
	checkErr(err)

	fmt.Printf("%d ID:Email pairs returned == %d ID:Email pairs expected\n", len(result[0]), len(inDB[0]))
	// checks that only values in DB are returned
	assert.True(t, len(result[0]) == len(inDB[0]))
}

func TestNewUnsubscribes(t *testing.T) {
	dao := databaseAccessObj.New()
	// Clear values in tables for clean test
	dao.Clear()

	// Increasing these values may produce false positives
	dataSet := bloomDataGenerator.GenData(1, 1000, 2000)
	pyld := Payload{0, dataSet[0]}

	data, err := json.Marshal(pyld)
	checkErr(err)

	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", bytes.NewBuffer(data))

	var result map[int][]string

	// Read the result
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	// converts the decoded result back to a Result struct
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error: Unable to unmarshal the request body. %v\n", err.Error())
	}

	fmt.Printf("%d ID:Email pairs returned == 0 ID:Email pairs expected\n", len(result))
	assert.True(t, len(result) == 0)

	// Insert the true values into the database
	dao.Insert(dataSet)

	// Update the bloom filter bit array
	_, err = http.Get(updateEndpoint)
	if err != nil {
		log.Printf("Error: Unable to update bit array. %v\n", err.Error())
	}

	res, _ = http.Post(membershipEndpoint, "application/json; charset=utf-8", bytes.NewBuffer(data))

	body, err = ioutil.ReadAll(res.Body)
	checkErr(err)

	var result2 map[int][]string
	// converts the decoded result back to a Result struct
	err = json.Unmarshal(body, &result2)
	checkErr(err)

	// Check again to see if bit array is updated
	fmt.Printf("%d ID:Email pairs returned == %d ID:Email pairs expected\n", len(result2[0]), len(dataSet[0]))
	assert.True(t, len(result2[0]) == len(dataSet[0]))
}

func TestResubscribed(t *testing.T) {
	dao := databaseAccessObj.New()

	// Clear tables for clean tests
	dao.Clear()

	var allData []string
	dataSet := bloomDataGenerator.GenData(1, 1000, 2000)
	extra := bloomDataGenerator.GenData(1, 1000, 2000)

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

	pyld := Payload{0, dataSum[0]}

	data, err := json.Marshal(pyld)
	checkErr(err)

	// get the memberships using yall the data
	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", bytes.NewBuffer(data))

	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	var result map[int][]string
	// converts the decoded result back to a Result struct
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Unable to unmarshal result body. %v", err.Error())
	}

	fmt.Printf("%d total ID:Email pairs != %d ID:Email pairs returned\n", len(dataSum[0]), len(result[0]))
	assert.False(t, len(dataSum[0]) == len(result[0]))
}
