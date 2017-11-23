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
	cmdPtr := flag.String("cmd", "", `
	Possible commands: 
	mktbls: Create 15 tables named unsub 0 - 14 in unsubscribed
		(no arguments needed) 
	repopulate: Clear all current values in tables and repopulates it 
		(can optionally specify number of users, mininum and maximum emails per user)
	add: Create new dataset and add it to the database
		(can optionally specify number of users, mininum and maximum emails per user)
	`)
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
	} else {
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
		default:
			fmt.Fprintf(os.Stderr, "Error: invalid command.\n")
			flag.PrintDefaults()
		}
	}
}
