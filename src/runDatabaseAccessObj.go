package main // probably need to convert this to a proper go test

import (
	"Inf191BloomFilter/src/bloomDataGenerator"
	"Inf191BloomFilter/src/databaseAccessObj"
	"fmt"
	"os"
	"strconv"
	"time"
)


func checkErr(err error){
	// check error from database if any
	if err != nil {
		panic(err)
	}
}


func main(){
	// command line inputs
	clInputs := os.Args[1:]

	// number of user ids
	numUsers, err := strconv.Atoi(clInputs[0])
	checkErr(err)

	// minimum and maximum number of emails per user_id
	minEmails, err := strconv.Atoi(clInputs[1])
	checkErr(err)

	maxEmails, err := strconv.Atoi(clInputs[2])
	checkErr(err)

	// log into db and clear table
	update := databaseAccessObj.New("bloom:test@/unsubscribed")
	fmt.Println("Clearing current db...")
	update.Clear()
	fmt.Println("Done.")

	// benchmarking for creating random data
	fmt.Printf("Generating test data (%d users, %d min addrs,  %d max addrs)...\n", numUsers, minEmails, maxEmails)
	start := time.Now()
	data := bloomDataGenerator.GenData(numUsers, minEmails, maxEmails)
	elapsed := time.Since(start)
	fmt.Printf("Done. Took %s\n", elapsed)

	// benchmarking for insert random data into one table in the db
	fmt.Println("Inserting test data into db...")
	start = time.Now()
	update.InsertDataSet(data)
	elapsed = time.Since(start)
	fmt.Printf("Done. Took %s\n\n", elapsed)

	/*
	fmt.Println("Clearing current db...")
	update.Clear()
	fmt.Println("Done.\n")

	// benchmarking for inserting random data into multiple shards
	fmt.Println("Inserting test data into db shards...")
	start = time.Now()
	update.InsertDataShards(data)
	elapsed = time.Since(start)
	fmt.Printf("Done. Took %s\n", elapsed)
	*/

	update.CloseConnection()
}
