package main

import (
	"fmt"
	// "databaseAccessObj"
	"net/http"
)

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
	fmt.Println("Testing...")
	fmt.Fprintf(w, "testing...")
}
