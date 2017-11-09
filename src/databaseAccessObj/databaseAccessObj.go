package databaseAccessObj

import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func checkErr(err error){
	if err != nil{
		panic(err)
	}
}

type Update struct {
	db *sql.DB
}

func New(dsn string) *Update{
	// input: dsn = 'username:password@/database'
	// create new connection
	db , err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return &Update{db}
}

func (update *Update)hasTable(databaseName, tableName string) bool{
	// For testing purposes
	// Checks if the specify tablename exist in the specified database
	db := update.db

	var check string
	err := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ? LIMIT 1", databaseName, tableName).Scan(&check)
	if(err == sql.ErrNoRows){
		return false
	} else{
		checkErr(err)
	}
	return true
}

func (update *Update)dropTable(tableName string){
	// for testing purposes
	// Delete the table in the db if exists
	db := update.db
	_, err := db.Exec("DROP TABLE IF EXISTS " + tableName)
	checkErr(err)
}

func (update *Update)EnsureTable(tableName string){
	// Checks if the specified table exists, if  not, Create a table
	// with that name

	db := update.db
	// Need to double check if this statment needs
	// placeholder
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ` + tableName +
	` (user_id int(11) NOT NULL DEFAULT '0',
	email varchar(255) NOT NULL DEFAULT '',
	PRIMARY KEY (user_id, email))
	ENGINE=InnoDB DEFAULT CHARSET=utf8`)
	checkErr(err)
}

func (update *Update)CloseConnection(){
	// Update method, closes connection
	update.db.Close()
}

func (update *Update)SelectAll() (map[int][]string){
	// Update method, select all rows from table: unsub_0 
	db := update.db
	result := map[int][]string{}
	rows, err := db.Query("SELECT user_id, email FROM unsub_0") // get all rows from database
	defer rows.Close()
	for rows.Next(){
		var user_id int
		var email string
		err = rows.Scan(&user_id, &email)
		checkErr(err)
		result[user_id] = append(result[user_id], email)
	}
	checkErr(err)
	return result
}

func (update *Update)InsertDataSet(dataSet map[int][]string){
	// Update method
	// Takes a (int, string[])map of data and insert them
	// into one table in the database
	db := update.db
	stmt, err := db.Prepare(`INSERT INTO unsub_0 (user_id, email) VALUES (?,?)`)
	checkErr(err)

	_, err = db.Exec("BEGIN")
	checkErr(err)

	for userid, emails := range dataSet{
		for i := range(emails){
			_, err := stmt.Exec(userid, emails[i])
			checkErr(err)
		}
	}

	_, err = db.Exec("COMMIT")
	checkErr(err)
}

func (update *Update)InsertDataShards(dataSet map[int][]string){
	// Takes a (int, string[])map of data and insert them
	// into different table according to the user_id
	// This should be used when trying to insert a larger amount
	// it would be faster to use InsertDataSet with smaller numbers
	db := update.db
	for userid, emails := range dataSet{
		tableName := "unsub_" + strconv.Itoa(userid)

		update.EnsureTable(tableName)
		stmt, err := db.Prepare(`INSERT INTO ` + tableName +
		` (user_id, email) VALUES (?,?)`)
		checkErr(err)

		_, err = db.Exec("BEGIN")
		checkErr(err)

		for i := range(emails){
			_, err := stmt.Exec(userid, emails[i])
			checkErr(err)
		}

		_, err = db.Exec("COMMIT")
		checkErr(err)
	}
}

func (update *Update)Clear(){
	// Delete all rows from a table in the database
	// Be very careful when using this function! It can
	// take a while to repopulate the db
	db := update.db
	_, err := db.Exec("TRUNCATE TABLE unsub_0")
	checkErr(err)
}

