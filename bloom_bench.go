package main

import (
	"Inf191BloomFilter/bloomDataGenerator"
	"Inf191BloomFilter/bloomManager"
	"Inf191BloomFilter/databaseAccessObj"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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

func getBFStats(bitArrSize, dbSize uint) float64 {
	numHashFunc := uint(10)
	bf := bloomManager.New(bitArrSize, numHashFunc)
	return bf.GetStats(dbSize)
}

func benchFalsePositiveProbility(dao *databaseAccessObj.Update, fromDBSize, toDBSize uint) {
	var prob float64
	bitArrSize := uint(100000)
	for fromDBSize < toDBSize || prob > 0.01 {
		prob = getBFStats(bitArrSize, fromDBSize)
		dao.LogTestResult("falposprob_bitarr_prob_"+strconv.Itoa(int(fromDBSize)), float64(bitArrSize), float64(prob))
		if prob > 0.01 {
			bitArrSize *= 2
		} else {
			fmt.Printf("Ran false-positive probility benchmark with database size of %d rows\n", int(fromDBSize))
			bitArrSize = 100000
			fromDBSize *= 2
		}
	}
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
	dao.ClearTestResults()
	benchFalsePositiveProbility(dao, 100000, 5000000)

	res := testing.Benchmark(BenchmarkUpdate1000)
	dao.LogTestResult("update_size_timeperop", float64(1000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate2000)
	dao.LogTestResult("update_size_timeperop", float64(2000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate4000)
	dao.LogTestResult("update_size_timeperop", float64(4000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate8000)
	dao.LogTestResult("update_size_timeperop", float64(8000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate16000)
	dao.LogTestResult("update_size_timeperop", float64(16000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUpdate32000)
	dao.LogTestResult("update_size_timeperop", float64(32000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership1000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(1000), float64(res.NsPerOp()))
	total_req := float64(res.N * 1000)
	total_time := res.T
	dao.LogTestResult("unsubmembership_size_reqpersec", float64(1000), total_req/total_time.Seconds())

	res = testing.Benchmark(BenchmarkUnsubMembership2000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(2000), float64(res.NsPerOp()))
	total_req = float64(res.N * 2000)
	total_time = res.T
	dao.LogTestResult("unsubmembership_size_reqpersec", float64(2000), total_req/total_time.Seconds())

	res = testing.Benchmark(BenchmarkUnsubMembership4000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(4000), float64(res.NsPerOp()))
	total_req = float64(res.N * 4000)
	total_time = res.T
	dao.LogTestResult("unsubmembership_size_reqpersec", float64(4000), total_req/total_time.Seconds())

	res = testing.Benchmark(BenchmarkUnsubMembership8000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(8000), float64(res.NsPerOp()))
	total_req = float64(res.N * 8000)
	total_time = res.T
	dao.LogTestResult("unsubmembership_size_reqpersec", float64(8000), total_req/total_time.Seconds())

	res = testing.Benchmark(BenchmarkUnsubMembership16000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(16000), float64(res.NsPerOp()))
	total_req = float64(res.N * 16000)
	total_time = res.T
	dao.LogTestResult("unsubmembership_size_reqpersec", float64(16000), total_req/total_time.Seconds())

	res = testing.Benchmark(BenchmarkUnsubMembership32000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(32000), float64(res.NsPerOp()))
	total_req = float64(res.N * 32000)
	total_time = res.T
	dao.LogTestResult("unsubmembership_size_reqpersec", float64(32000), total_req/total_time.Seconds())

	dao.CloseConnection()
}
