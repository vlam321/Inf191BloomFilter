package main

import "fmt"
import "github.com/marpaia/graphite-golang"

import "Inf191BloomFilter/databaseAccessObj"

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	g, err := graphite.NewGraphite("localhost", 2003)
	checkErr(err)

	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	results := dao.SelectTestResults()
	for i := 0; i < len(results); i++ {
		err = g.SimpleSend("Test_Results_1", fmt.Sprintf("%f", results[i].X))
		err = g.SimpleSend("Test_Results_2", fmt.Sprintf("%f", results[i].Y))
		checkErr(err)
	}

	/*
		if err != nil {
			panic(err)
		}
		for i := 0; i < 20; i++ {
			time.Sleep(time.Second)
			err = g.SimpleSend("foo.bar.metric_one", fmt.Sprintf("%d", i))
			if err != nil {
				log.Printf("err: %v", err)
				return
			}
			err = g.SimpleSend("foo.bar.metric_two", fmt.Sprintf("%d", 20-i))
			if err != nil {
				log.Printf("err: %v", err)
				return
			}
		}
		fmt.Println("hello world")
	*/
}
