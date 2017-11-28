/* This file will contain functionalities to host the main service
in which the client will make requests to. Use net/http to create
API endpoints to access the following functionalities:
	- UpdateBloomFilter
		- AddToBloomFilter
		- RepopulateBloomFilter
	- GetArrayOfUnsubscribedEmails
*/
package main

import (
	"Inf191BloomFilter/bloomManager"
	"encoding/json"
	"fmt"
	"net/http"
)

type payload struct {
	UserId int
	Emails []string
}

//Global variable
//The bloom filter for this server
var bf = bloomManager.New()


//handleUpdate will update the respective bloomFilter
//server will keep track of when the last updated time is. to call 
//update every _ time. 
// func handleUpdate(r http.ResponseWriter, req *http.Request) {
// 	if lastUpdate.Sub(time.Now()).Seconds == time.Minute{

// 	}
// 	bf.UpdateBloomFilter()

// }

//checkErr checks fro errors in the decoded text and encoded text... 
func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//handleFilterUnsubscripe handles when the client requests to recieve 
//those emails that are unsubscribed therefore, IN the database. 
func handleFilterUnsubscribed(w http.ResponseWriter, r *http.Request) {
	
	var buff []byte
	var payload payload

	//Result struct made to carry the result of unsuscribed emails 
	type Result struct {
	Trues []string
	}

	//decode byets from request body
	err := json.NewDecoder(r.Body).Decode(&buff)

	//check for error; if error write 404
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
		// http.Error(w, err.Error(), 400)
		// return
	}
	
	// converts the decoded result to a payload struct
	err = json.Unmarshal(buff, &payload)
	checkErr(err)

	fmt.Println(payload)

	//uses bloomManager to get the result of unsubscribed emails 
	//puts them in struct, result 
	emails := bf.GetArrayOfUnsubscribedEmails(payload.Emails)
	filteredEmails := Result{emails}

	//convert result to json 
	js, err := json.Marshal(filteredEmails)
	checkErr(err)

	fmt.Println(js)

	//write back to client
	w.Write(js)
}


func timeManager(){
	ticker := time.NewTicker(time.Second * 2)	
	go func() {
        	for t := range ticker.C {
				//Call update bloom filter
				bf.UpdateBloomFilter()
            	fmt.Println("Bloom Filter Updated at: ", t.Format("2006-01-02 3:4:5 PM"))
	
       		}
	}()
	//Figure out how to run without sleep? 
	//use go forever? an option
	time.Sleep(time.Second * 10)
    	ticker.Stop()
    	fmt.Println("Ticker stopped")
}


func main() {
   
	timeManager()
	http.HandleFunc("/filterUnsubscribed", handleFilterUnsubscribed)
	// http.HandleFunc("/update", handleUpdate)
	
	http.ListenAndServe(":9090", nil)
}
