package main

import (
	"Inf191BloomFilter/bloomDataGenerator"
	"Inf191BloomFilter/databaseAccessObj"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

const dsn = "bloom:test@/unsubscribed"
const membershipEndpoint = "http://localhost:9090/filterUnsubscribed"
const updateEndpoint = "http://localhost:9090/update"

type Payload struct {
	UserId int
	Emails []string
}

type Result struct {
	Trues []string
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func repopulateDatabase(numValues int) {
	data := bloomDataGenerator.GenData(1, numValues, numValues+1)
	dao := databaseAccessObj.New(dsn)
	dao.Clear()
	dao.Insert(data)
	dao.CloseConnection()
}

func updateBitArray() {
	_, err := http.Get(updateEndpoint)
	checkErr(err)
}

func conv2Json(payload Payload) []byte {
	data, err := json.Marshal(payload)
	checkErr(err)
	return data
}

func conv2Buff(idEmailJson []byte) io.Reader {
	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(&idEmailJson)
	checkErr(err)
	return buff
}

func getUnsub(dataset map[int][]string) []string {
	idEmailpayload := Payload{0, dataset[0]}
	idEmailJson := conv2Json(idEmailpayload)
	buff := conv2Buff(idEmailJson)

	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", buff)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	var result Result
	err = json.Unmarshal(body, &result)
	checkErr(err)
	return result.Trues
}

func benchmarkBitArrayUpdate(numValues int, b *testing.B) {
	repopulateDatabase(numValues)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		updateBitArray()
	}
	b.StopTimer()
	fmt.Printf("Ran bit array update benchmark with %d values (%d iterations).\n", numValues, b.N)
}

func benchmarkUnsubMembership(numIDEmailPairs int, b *testing.B) {
	dao := databaseAccessObj.New(dsn)
	dataset := dao.SelectRandSubset(0, numIDEmailPairs)
	dao.CloseConnection()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = getUnsub(dataset)
	}
	b.StopTimer()
	fmt.Printf("Ran unsubscribed members filter benchmark with %d values (%d iterations).\n", numIDEmailPairs, b.N)
}

func BenchmarkUpdate1000(b *testing.B)  { benchmarkBitArrayUpdate(1000, b) }
func BenchmarkUpdate2000(b *testing.B)  { benchmarkBitArrayUpdate(2000, b) }
func BenchmarkUpdate4000(b *testing.B)  { benchmarkBitArrayUpdate(4000, b) }
func BenchmarkUpdate8000(b *testing.B)  { benchmarkBitArrayUpdate(8000, b) }
func BenchmarkUpdate16000(b *testing.B) { benchmarkBitArrayUpdate(16000, b) }
func BenchmarkUpdate32000(b *testing.B) { benchmarkBitArrayUpdate(32000, b) }

func BenchmarkUnsubMembership1000(b *testing.B)  { benchmarkUnsubMembership(1000, b) }
func BenchmarkUnsubMembership2000(b *testing.B)  { benchmarkUnsubMembership(2000, b) }
func BenchmarkUnsubMembership4000(b *testing.B)  { benchmarkUnsubMembership(4000, b) }
func BenchmarkUnsubMembership8000(b *testing.B)  { benchmarkUnsubMembership(8000, b) }
func BenchmarkUnsubMembership16000(b *testing.B) { benchmarkUnsubMembership(16000, b) }
func BenchmarkUnsubMembership32000(b *testing.B) { benchmarkUnsubMembership(32000, b) }

func main() {
	dao := databaseAccessObj.New(dsn)

	res := testing.Benchmark(BenchmarkUpdate1000)
	dao.LogTestResult("update_size_timeperop", float32(1000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate2000)
	dao.LogTestResult("update_size_timeperop", float32(2000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate4000)
	dao.LogTestResult("update_size_timeperop", float32(4000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate8000)
	dao.LogTestResult("update_size_timeperop", float32(8000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate16000)
	dao.LogTestResult("update_size_timeperop", float32(16000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate32000)
	dao.LogTestResult("update_size_timeperop", float32(32000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership1000)
	dao.LogTestResult("unsubmembership_size_timeperop", float32(1000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership2000)
	dao.LogTestResult("unsubmembership_size_timeperop", float32(2000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership4000)
	dao.LogTestResult("unsubmembership_size_timeperop", float32(4000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership8000)
	dao.LogTestResult("unsubmembership_size_timeperop", float32(8000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership16000)
	dao.LogTestResult("unsubmembership_size_timeperop", float32(18000), float32(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership32000)
	dao.LogTestResult("unsubmembership_size_timeperop", float32(32000), float32(res.NsPerOp()))

	dao.CloseConnection()
}
