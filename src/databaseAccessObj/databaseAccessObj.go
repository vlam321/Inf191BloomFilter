package databaseAccessObj

import(
	//"../bloomDataGenerator"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
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

func (update *Update)HasTable(databaseName, tableName string) bool{
	// Checks if the specify tablename exist in the specified database
	db := update.db
	_, err := db.Exec("use information_schema")
	row, err := db.Query(`SELECT count(table_name) FROM tables
	WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`, databaseName, tableName)
	checkErr(err)
	var count int
	for row.Next(){
		err := row.Scan(&count)
		checkErr(err)
	}
	if (count > 0){return true}
	return false
}

func (update *Update)CloseConnection(){
	// Update method, closes connection
	update.db.Close()
}

func (update *Update)SelectAll() (map[int][]string){
	// select all from database  
	db := update.db
	result := map[int][]string{}
	fmt.Println(result)
	rows, err := db.Query("SELECT user_id, email FROM unsub_0") // get all rows from database
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
	// Takes a (int, string[])map of data and insert them
	// into the specified db
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

func (update *Update)Clear(){
	// Delete all rows from a table in the database
	// Be very careful when using this function! It can
	// take a while to repopulate the db
	db := update.db
	_, err := db.Exec("TRUNCATE TABLE unsub_0")
	checkErr(err)
}

