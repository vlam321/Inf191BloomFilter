package main // probably need to convert this to a proper go test

import (
	"Inf191BloomFilter/bloomDataGenerator"
	"Inf191BloomFilter/databaseAccessObj"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type UserInputs struct {
	command  string
	numUser  int
	minEmail int
	maxEmail int
}

const dsn = "bloom:test@/unsubscribed"
const schema = `(user_id int(11), email varchar(255), ts timestamp default current_timestamp, primary key (user_id, email));`

func checkErr(err error) {
	// check error from database if any
	if err != nil {
		panic(err)
	}
}

// Grabs the command line argments
// and return a Inputs object containinf the
// values. And a nil if no cli atgements were
// given
func getCommandLineInputs() UserInputs {
	cmdPtr := flag.String("cmd", "", "Possible commands: 'repopulate', 'add', 'del', 'mktbls'")
	userPtr := flag.Int("numUser", 1, "Possible inputs: integers > 0")
	minEmailPtr := flag.Int("minEmail", 1, "Possible inputs: integer > 0")
	maxEmailPtr := flag.Int("maxEmail", 2, "Possible inputs: integer > minEmail")
	flag.Parse()
	return UserInputs{*cmdPtr, *userPtr, *minEmailPtr, *maxEmailPtr}
}

// Given the user inputs, clear existing data and repopulate
// the table with new randomly generated data
func handleRepopulate(numUser, minEmail, maxEmail int) {
	dao := databaseAccessObj.New(dsn)
	dao.Clear()
	dataset := bloomDataGenerator.GenData(numUser, minEmail, maxEmail)
	dao.Insert(dataset)
	dao.CloseConnection()
}

func handleAdd(numUser, minEmail, maxEmail int) {
	dao := databaseAccessObj.New(dsn)
	dataset := bloomDataGenerator.GenData(numUser, minEmail, maxEmail)
	dao.Insert(dataset)
	dao.CloseConnection()
}

func handleDel(numUser, minEmail, maxEmail int) {
	dao := databaseAccessObj.New(dsn)
	// change arguments
	// use a int to determines how many rows to delete
	// need dao to return a subset of random useremail pairs
	// delete those from the db (may need a func in dao for that)
	// return the ones deleted?
	dao.CloseConnection()
}

func handleMkTbl() {
	dao := databaseAccessObj.New(dsn)
	for i := 0; i < 15; i++ {
		tablename := "unsub_" + strconv.Itoa(i)
		dao.MkTbl(tablename, schema)
	}
	dao.CloseConnection()
}

func main() {
	userInputs := getCommandLineInputs()
	if userInputs.command == "" {
		fmt.Fprintf(os.Stderr, "Error: cmd cannot be empty.\n")
		flag.PrintDefaults()
	}
	switch userInputs.command {
	case "mktbls":
		handleMkTbl()
		fmt.Println("Done. Created tables in unsubscribed.")
	case "repopulate":
		handleRepopulate(userInputs.numUser, userInputs.minEmail, userInputs.maxEmail)
		fmt.Printf("Done. \n")
	case "add":
		handleAdd(userInputs.numUser, userInputs.minEmail, userInputs.maxEmail)
		fmt.Printf("Done. \n")
	case "del":
		handleDel(userInputs.numUser, userInputs.minEmail, userInputs.maxEmail)
		fmt.Printf("Done. \n")
	default:
		fmt.Fprintf(os.Stderr, "Error: invalid command.\n")
		flag.PrintDefaults()
	}
}
