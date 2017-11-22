/* This file will contain functionality for a simulated
client. The client is responsible for makine http requests
to the dbServer and bloomFilterServer
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const membershipEndpoint = "http://localhost:9090/filterUnsubscribed"

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func TestUnsub(t *testing.T) {
	payload := Payload{1, []string{"sodfd", "fdsafasd"}}
	buff := new(bytes.Buffer)

	data, err := json.Marshal(payload)
	checkErr(err)

	err = json.NewEncoder(buff).Encode(data)
	checkErr(err)

	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", buff)
	fmt.Println(res)
}
