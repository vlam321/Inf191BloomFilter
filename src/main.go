package main

// import "./bloomDataGenerator"
import "fmt"
import "net/http"
import "strings"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

func getAverageLen(addrs []string)(int){
	numEle := len(addrs)
	numChar := 0
	for i := range addrs{
		numChar += len([]rune(addrs[i]))
	}
	return numChar/numEle
}

func main(){
	// Connecting to the DB
	db, err := sql.Open("mysql", "username:password@/unsubscribed")
	if err != nil {     panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic 
	}

	// Closing the connection at the end
	defer db.Close()

	// testing the connection
	err = db.Ping()
	if err != nil {     panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic 
	}

	// Populating the DB (run once)
	// _, err = db.Exec("insert into unsub_0 (user_id, email) values (180, ?)", "vlam321@gamil.com")

	// Querying the DB
	rows, err := db.Query("SELECT user_id, email FROM unsub_0;")
	if err != nil {     panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic 
	}

	// Printing the results
	for rows.Next(){
		var userid int
		var email string
		err = rows.Scan(&userid, &email)
		if err != nil{
			panic(err.Error())
		}
		fmt.Printf("%d, %s\n", userid, email)
	}

	// Check for errors in rows
	err = rows.Err()
	if err != nil{
		panic(err.Error())
	}
	fmt.Println("Success!")

	// net/http stuffs
	http.HandleFunc("/hello", sayhelloName) // set router
	err = http.ListenAndServe(":9090", nil) // set listen port 
	if err != nil { panic(err.Error()) }
}

// func used for http stuffs
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  // parse arguments, you have to call this by yourself
	fmt.Println(r.Form)  // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // send data to client side 
}

