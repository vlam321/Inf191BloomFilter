package main

import (
	"fmt"
	"Inf191BloomFilter/src/databaseAccessObj"
	"net/http"
	"encoding/json"
)

const dsn = "bloom:test@/unsubscribed"

func checkErr(err error){
	if err != nil {
		panic(err.Error())
	}
}

func main(){
	http.HandleFunc("/insertUserEmail", insertUserEmail)
	err := http.ListenAndServe(":9090", nil)
	checkErr(err)
}

func insertUserEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inserting new data...")
	update := databaseAccessObj.New(dsn)
	var data map[int][]string
	if r.Body == nil {
		http.Error(w, "No data received", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil{
		http.Error(w, err.Error(), 400)
		return
	}
	update.InsertDataSet(data)
	fmt.Println("Done .")
}

