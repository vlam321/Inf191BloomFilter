/* This file will contain functionality for a simulated
client. The client is responsible for makine http requests
to the dbServer and bloomFilterServer
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const membershipEndpoint = "http://localhost:9090/filterUnsubscribed"

func TestUnsub(t *testing.T) {
	payload := Payload{1, []string{"sodfd", "fdsafasd"}}
	buff := new(bytes.Buffer)

	data, err := json.Marshal(payload)
	checkErr(err)

	err = json.NewEncoder(buff).Encode(data)
	checkErr(err)

	var deres []int8
	var arr Result
	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", buff)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(body)

	err = json.NewDecoder(res.Body).Decode(&deres)

	err = json.Unmarshal(body, &arr)
	checkErr(err)
	fmt.Println(arr.Trues)

	// var buff2 []byte
	// var payload2 Payload

	// err2 := json.NewDecoder(res.Body).Decode(&buff2)
	// checkErr(err2)
	// err2 = json.Unmarshal(buff2, &payload2)
	// checkErr(err2)

	// fmt.Println(payload2)

}
