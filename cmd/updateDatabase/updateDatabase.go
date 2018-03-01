package main // probably need to convert this to a proper go test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/vlam321/Inf191BloomFilter/bloomDataGenerator"
	"github.com/vlam321/Inf191BloomFilter/databaseAccessObj"
)

type UserInputs struct {
	command   string
	numUser   int
	minEmail  int
	maxEmail  int
	tableNum  int
	numEmail  int
	avgEmails int
	interval  int
}

const unsub_schema = `(user_id int(11), email varchar(255), ts timestamp default current_timestamp, primary key (user_id, email));`
const test_result_schema = `(result_type varchar(30) NOT NULL, x_axis float NOT NULL, y_axis float NOT NULL );`

// getCommandLineInputs returns object of user input; nil if no input
func getCommandLineInputs() UserInputs {
	cmdPtr := flag.String("cmd", "", `
	Possible commands: mktbls, repopulate, add, delete
	mktbls: Create 15 tables named unsub 0 - 14 in unsubscribed
		(no arguments needed) 
	repopulate: Clear all current values in tables and repopulates it 
		(can optionally specify number of users, mininum and maximum emails per user)
	add: Create new dataset and add it to the database
		(can optionally specify number of users, mininum and maximum emails per user)
	del: Delete dataset from database
	auto: automatically insert data into database
		(can optionally specify number of emails to insert and interval between inserts (in ms))
	`)
	userPtr := flag.Int("numUser", 1, "Possible inputs: integers > 0")
	minEmailPtr := flag.Int("minEmail", 1, "Possible inputs: integer > 0")
	maxEmailPtr := flag.Int("maxEmail", 2, "Possible inputs: integer > minEmail")
	tableNumPtr := flag.Int("tableNum", 0, "Possible inputs: 0-14")
	numEmailPtr := flag.Int("numEmail", 0, "Possible inputs: >= 0")
	avgEmailsPtr := flag.Int("avgEmails", 1, "Possible inputs: >= 0")
	intervalPtr := flag.Int("interval", 1, "Possible inputs:")
	flag.Parse()
	return UserInputs{*cmdPtr, *userPtr, *minEmailPtr, *maxEmailPtr, *tableNumPtr, *numEmailPtr, *avgEmailsPtr, *intervalPtr}
}

// handleRepopulate clears database and populates with random data based on input
func handleRepopulate(numUser, minEmail, maxEmail int) {
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	dao.Clear()
	dataset := bloomDataGenerator.GenData(numUser, minEmail, maxEmail)
	dao.Insert(dataset)
}

// handleAdd adds random data based on input to db
func handleAdd(numUser, minEmail, maxEmail int) {
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	dataset := bloomDataGenerator.GenData(numUser, minEmail, maxEmail)
	dao.Insert(dataset)
}

// handleDelete takes a table number and a number of rows and remove them from the db
func handleDel(tableNum, numEmail int) {
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	dao.Delete(tableNum, numEmail)
}

// handleMakeTable creates all tables necessary in db
func handleMakeTable() {
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	for i := 0; i < 15; i++ {
		tablename := "unsub_" + strconv.Itoa(i)
		dao.MakeTable(tablename, unsub_schema)
	}
	dao.MakeTable("test_results", test_result_schema)
}

// autoUpdate updates database based on time input
func handleAuto(avgEmails int, interval int) {
	dao := databaseAccessObj.New()
	defer dao.CloseConnection()
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for _ = range ticker.C {
		log.Printf("Inserting %d emails into database", avgEmails)
		dao.Insert(bloomDataGenerator.GenData(1, avgEmails, avgEmails+1))
	}
}

func main() {
	userInputs := getCommandLineInputs()
	if userInputs.command == "" {
		fmt.Fprintf(os.Stderr, "Error: cmd cannot be empty.\n")
		flag.PrintDefaults()
	} else {
		switch userInputs.command {
		case "mktbls":
			handleMakeTable()
			fmt.Printf("Done. Created tables in unsubscribed.\n")
		case "repopulate":
			handleRepopulate(userInputs.numUser, userInputs.minEmail, userInputs.maxEmail)
			fmt.Printf("Done. \n")
		case "add":
			handleAdd(userInputs.numUser, userInputs.minEmail, userInputs.maxEmail)
			fmt.Printf("Done. \n")
		case "del":
			handleDel(userInputs.tableNum, userInputs.numEmail)
			fmt.Printf("Done. \n")
		case "auto":
			handleAuto(userInputs.avgEmails, userInputs.interval)
		default:
			fmt.Fprintf(os.Stderr, "Error: invalid command.\n")
			flag.PrintDefaults()
		}
	}
}
