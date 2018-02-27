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
	"time"

	"strconv"

	"github.com/vlam321/Inf191BloomFilter/bloomDataGenerator"
	"github.com/vlam321/Inf191BloomFilter/databaseAccessObj"
	"github.com/vlam321/Inf191BloomFilter/payload"
)

const dockerEp = "http://192.168.99.100:9090/filterUnsubscribed"
const awsEp = "http://54.183.120.79:9090/filterUnsubscribed"
const shard1 = "http://13.56.155.181:9090/filterUnsubscribed"

var numShards int

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
		result[temp.UserId] = temp.Emails
	}
	return result
}

// skip marshal/unmarshal
/*
func getUnsub(pl [][]byte) [][]byte{
	var result [][]byte
	for i:=0; i<len(pl); i++{
		res, err := http.Post(awsEp, "application/json; charset=utf-8", bytes.NewBuffer(pl[i]))
		if err != nil{
			log.Printf("Error processing request: %v\n", err)
			return nil
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil{
			log.Printf("Error reading body: %v\n", err)
			return nil
		}
		result = append(result, body)
	}
	return result
}
*/

func benchmarkUnsub(trueResults, falseResults int, b *testing.B) {
	falseResults = 0
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	testData := make(map[int][]string)

	// skip marshal/unmarshal
	/*
		var pl []payload.Payload
		var plJson [][]byte
	*/

	for i := 0; i < numShards; i++ {
		testData[i] = append(dao.SelectRandSubset(i, trueResults)[i], bloomDataGenerator.GenData(1, falseResults, falseResults+1)[0]...)
		// skip marshal/unmarshal
		/*
			pl = append(pl, payload.Payload{i, testData[i]})
			tempJson, err := json.Marshal(pl[i])
			if err != nil{
				log.Printf("error marshaling: %v\n", err)
				return
			}
			plJson = append(plJson, tempJson)
		*/
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getUnsub(testData)
		//getUnsub(plJson)
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

func BenchmarkUnsub1000(b *testing.B)    { benchmarkUnsub(1000, 1000, b) }
func BenchmarkUnsub2000(b *testing.B)    { benchmarkUnsub(2000, 2000, b) }
func BenchmarkUnsub4000(b *testing.B)    { benchmarkUnsub(4000, 4000, b) }
func BenchmarkUnsub8000(b *testing.B)    { benchmarkUnsub(8000, 8000, b) }
func BenchmarkUnsub16000(b *testing.B)   { benchmarkUnsub(16000, 16000, b) }
func BenchmarkUnsub50000(b *testing.B)   { benchmarkUnsub(50000, 50000, b) }
func BenchmarkUnsub100000(b *testing.B)  { benchmarkUnsub(100000, 100000, b) }
func BenchmarkUnsub500000(b *testing.B)  { benchmarkUnsub(500000, 500000, b) }
func BenchmarkUnsub1000000(b *testing.B) { benchmarkUnsub(1000000, 1000000, b) }

func main() {
	// go run cmd/bloomBench/bloomRouterBench/bloomRouterBench.go [number of working db shards]
	log.Printf("starting bloom router bench\n")
	numShards, _ = strconv.Atoi(os.Args[1])
	f, err := os.Create(fmt.Sprintf("log/bench%d_%s.txt", numShards, time.Now().Format("2006-01-02T15-04-05AM")))
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
	res = testing.Benchmark(BenchmarkUnsub500000)
	log.Println(res.NsPerOp())
	logResult(f, 500000, res.NsPerOp())
	res = testing.Benchmark(BenchmarkUnsub1000000)
	log.Println(res.NsPerOp())
	logResult(f, 1000000, res.NsPerOp())
}
