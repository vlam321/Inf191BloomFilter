package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/vlam321/Inf191BloomFilter/bloomDataGenerator"
	"github.com/vlam321/Inf191BloomFilter/databaseAccessObj"
	"github.com/vlam321/Inf191BloomFilter/payload"
)

const dockerEp = "http://192.168.99.100:9090/filterUnsubscribed"
const awsEp = "http://13.56.59.216:9090/filterUnsubscribed"
const numShards = 2

func getUnsub(dataset map[int][]string) map[int][]string {
	result := make(map[int][]string)
	for k, v := range dataset {
		pl := payload.Payload{k, v}
		plJson, err := json.Marshal(pl)
		if err != nil {
			log.Printf("Error json marshaling: %v\n", err)
			return nil
		}
		res, _ := http.Post(awsEp, "application/json; charset=utf-8", bytes.NewBuffer(plJson))
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading body: %v\n", err)
			return nil
		}
		var temp payload.Payload
		err = json.Unmarshal(body, &temp)
		if err != nil {
			log.Printf("Error unmarshaling body: %v\n", err)
			return nil
		}
	}
	return result
}

func benchmarkUnsub(trueResults, falseResults int, b *testing.B) {
	falseResults = 0
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	testData := make(map[int][]string)
	for i := 0; i < numShards; i++ {
		testData[i] = append(dao.SelectRandSubset(i, trueResults)[i], bloomDataGenerator.GenData(1, falseResults, falseResults+1)[0]...)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getUnsub(testData)
	}
	b.StopTimer()
	log.Printf("Ran unsubscribed members filter benchmark: %d shards %d true values and %d false values (%d iterations).\n", numShards, trueResults, falseResults, b.N)
}

func logResult(f *os.File, n int, result int64) {
	_, err := f.WriteString(fmt.Sprintf("%d %d\n", n, result))
	if err != nil {
		log.Printf("Error writing file: %v\n", err)
		return
	}

}

func BenchmarkUnsub1000(b *testing.B)   { benchmarkUnsub(1000, 1000, b) }
func BenchmarkUnsub2000(b *testing.B)   { benchmarkUnsub(2000, 2000, b) }
func BenchmarkUnsub4000(b *testing.B)   { benchmarkUnsub(4000, 4000, b) }
func BenchmarkUnsub8000(b *testing.B)   { benchmarkUnsub(8000, 8000, b) }
func BenchmarkUnsub16000(b *testing.B)  { benchmarkUnsub(16000, 16000, b) }
func BenchmarkUnsub50000(b *testing.B)  { benchmarkUnsub(50000, 50000, b) }
func BenchmarkUnsub100000(b *testing.B) { benchmarkUnsub(100000, 100000, b) }

func main() {
	f, err := os.Create(fmt.Sprintf("log/bench%d.txt", numShards))
	if err != nil {
		log.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("numShards: %d\n", numShards))
	if err != nil {
		log.Printf("Error writing numShards: %v\n", err)
		return
	}

	res := testing.Benchmark(BenchmarkUnsub1000)
	log.Println(res.NsPerOp())
	logResult(f, 1000, res.NsPerOp())
	res = testing.Benchmark(BenchmarkUnsub2000)
	log.Println(res.NsPerOp())
	logResult(f, 2000, res.NsPerOp())
	res = testing.Benchmark(BenchmarkUnsub4000)
	log.Println(res.NsPerOp())
	logResult(f, 4000, res.NsPerOp())
	res = testing.Benchmark(BenchmarkUnsub8000)
	log.Println(res.NsPerOp())
	logResult(f, 8000, res.NsPerOp())
	res = testing.Benchmark(BenchmarkUnsub16000)
	log.Println(res.NsPerOp())
	logResult(f, 16000, res.NsPerOp())
	res = testing.Benchmark(BenchmarkUnsub50000)
	log.Println(res.NsPerOp())
	logResult(f, 50000, res.NsPerOp())
	res = testing.Benchmark(BenchmarkUnsub100000)
	log.Println(res.NsPerOp())
	logResult(f, 100000, res.NsPerOp())
}
