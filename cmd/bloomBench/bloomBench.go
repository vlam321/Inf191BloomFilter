package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"

	"Inf191BloomFilter/bloomDataGenerator"
	"Inf191BloomFilter/bloomManager"
	"Inf191BloomFilter/databaseAccessObj"
)

const membershipEndpoint = "http://localhost:9090/filterUnsubscribed"
const updateEndpoint = "http://localhost:9090/update"

type Payload struct {
	UserId int
	Emails []string
}

type Result struct {
	Trues []string
}

// repopulateDatabase clears db then adds random data based on numValues
func repopulateDatabase(numValues int) {
	data := bloomDataGenerator.GenData(1, numValues, numValues+1)
	dao := databaseAccessObj.New()
	dao.Clear()
	dao.Insert(data)
	dao.CloseConnection()
}

// updateBitArray makes request to updateEndpoint to update bf
func updateBitArray() {
	_, err := http.Get(updateEndpoint)
	if err != nil {
		log.Printf("Error updating bit array: %v\n", err)
		return
	}
}

// conv2Json converts payload input into JSON
func conv2Json(payload Payload) []byte {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error json marshaling: %v\n", err)
		return nil
	}
	return data
}

// getBFStats returns false positive rate of bloom filter
func getBFStats(bitArrSize, dbSize uint) float64 {
	numHashFunc := uint(10)
	bf := bloomManager.New(bitArrSize, numHashFunc)
	return bf.GetStats(dbSize)
}

// benchFalsePositiveProbability benchmarks different false postive rates based on db size
func benchFalsePositiveProbability(dao *databaseAccessObj.Conn, fromDBSize, toDBSize uint) {
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

// getUnsub
func getUnsub(dataset map[int][]string) map[int][]string {
	idEmailpayload := Payload{0, dataset[0]}
	idEmailJson := conv2Json(idEmailpayload)

	res, _ := http.Post(membershipEndpoint, "application/json; charset=utf-8", bytes.NewBuffer(idEmailJson))

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		return nil
	}

	var result map[int][]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error unmarshaling body: %v\n", err)
		return nil
	}
	return result
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
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	dataset := dao.SelectRandSubset(0, numIDEmailPairs)
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
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	dao.ClearTestResults()
	benchFalsePositiveProbability(dao, 100000, 5000000)

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

	res = testing.Benchmark(BenchmarkUnsubMembership2000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(2000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership4000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(4000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership8000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(8000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership16000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(16000), float64(res.NsPerOp()))

	res = testing.Benchmark(BenchmarkUnsubMembership32000)
	dao.LogTestResult("unsubmembership_size_timeperop", float64(32000), float64(res.NsPerOp()))

}
