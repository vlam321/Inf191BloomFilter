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

func TestUnsub(t *testing.T){
	payload := Payload{1, []string{"sodfd", "fdsafasd"}}
	buff := new(bytes.Buffer)

	data, err := json.Marshal(payload)
	checkErr(err)

	err = json.NewEncoder(buff).Encode(data)
	checkErr(err)

	res, _ := http.Post(membership_endpoint, "application/json; charset=utf-8", buff)
	fmt.Println(res)
}
