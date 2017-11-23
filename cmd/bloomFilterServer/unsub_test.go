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
	"io/ioutil"
	"net/http"
	"testing"
)

const membershipEndpoint = "http://localhost:9090/filterUnsubscribed"

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
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	var dataSum []string
	buff := new(bytes.Buffer)
	var payload Payload

	// Generate random id_email pairs (positives) and save it in a var
	inDB := bloomDataGenerator.GenData(1, 100, 200)

	// Insert new data in the db
	dao.Insert(inDB)

	// Call BF server to update the bit array
	/* Wait for service to be completed
	res, err := http.Get(update_bit_array_endpoint)
	checkErr(err)
	*/

	// Generate more raandom id_email pairs (negatives) and save it ina var
	notInDB := bloomDataGenerator.GenData(1, 50, 100)

	// Concatenate the two slices
	for userid, emails := range inDB {
		dataSum = append(emails, notInDB[userid]...)
	}

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

	var deres []int8
	var arr Result

	// Read the result
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	// decode the body to []int8 to be unmarshaled
	err = json.NewDecoder(res.Body).Decode(&deres)

	// converts the decoded result back to a Result struct
	err = json.Unmarshal(body, &arr)
	checkErr(err)

	// checks that only values in DB are returned
	/* Wait for service to be completed
	assert.True(t, len(arr.Trues) == len(inDB[0]))
	*/
}
