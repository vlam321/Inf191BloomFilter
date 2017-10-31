package main

import (
	"Inf191BloomFilter/src/bloomDataGenerator"
	"Inf191BloomFilter/src/updateBloomData"
	"fmt"
	"os"
	"strconv"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	clInputs := os.Args[1:]

	numUsers, err := strconv.Atoi(clInputs[0])
	checkErr(err)

	minEmails, err := strconv.Atoi(clInputs[1])
	checkErr(err)

	maxEmails, err := strconv.Atoi(clInputs[2])
	checkErr(err)

	update := updateBloomData.New("bloom:test@/unsubscribed")
	fmt.Println("Clearing current db...")
	update.Clear()
	fmt.Println("Done.")

	fmt.Printf("Generating test data (%d users, %d min addrs,  %d max addrs)...\n", numUsers, minEmails, maxEmails)
	start := time.Now()
	data := bloomDataGenerator.GenData(numUsers, minEmails, maxEmails)
	elapsed := time.Since(start)
	fmt.Printf("Done. Took %s\n", elapsed)

	fmt.Println("Inserting test data into db...")
	start = time.Now()
	update.InsertDataSet(data)
	elapsed = time.Since(start)
	fmt.Printf("Done. Took %s\n", elapsed)

	update.CloseConnection()
}
