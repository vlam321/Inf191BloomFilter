/* This file will contain functionality for a simulated
client. The client is responsible for makine http requests
to the dbServer and bloomFilterServer
*/

package main
import (
	"Inf191BloomFilter/bloomDataGenerator"
	"encoding/json"
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

const host = "http://localhost:9090/insertUserEmail"

func TestUnsub(t *testing.T){
	val := bloomDataGenerator.GenData(1, 2, 10)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(val)

	res, _ := http.Post(host, "application/json; charset=utf-8", b)
	fmt.Println(res)
}
