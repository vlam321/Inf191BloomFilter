package databaseAccessObj

import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"math"
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

func (update *Update)CloseConnection(){
	// Update method, closes connection
	update.db.Close()
}

func (update *Update)Select(dataSet map[int][]string) (map[int][]string){
	db := update.db
	result := make(map[int][]string)
	for userid, emails := range dataSet{
		tableName := "unsub_" + strconv.Itoa(int(math.Mod(float64(userid),15.0)))
		sqlStr := "SELECT user_id, email FROM " + tableName + " WHERE user_id = " + strconv.Itoa(userid) + " AND ("
		var vals []interface{}
		for i := range(emails){
			sqlStr += "email = ? OR "
			vals = append(vals, dataSet[userid][i])
		}
		sqlStr = sqlStr[0:len(sqlStr)-4]
		sqlStr += ")"
		rows, err := db.Query(sqlStr, vals...)
		if(err == sql.ErrNoRows){
			continue
		}
		defer rows.Close()
		for rows.Next(){
			var user_id int
			var email string
			err = rows.Scan(&user_id, &email)
			checkErr(err)
			result[user_id] = append(result[user_id], email)
		}
	}
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

func (update *Update)BatchInsert(dataSet map[int][]string){
	// Takes a (int, string[])map of data and insert them
	// into different table according to the user_id
	db := update.db
	for userid, emails := range dataSet{
		tableName := "unsub_" + strconv.Itoa(int(math.Mod(float64(userid),15.0)))
		sqlStr := "INSERT INTO " + tableName + "(user_id, email) VALUES "
		var vals []interface{}
		for i := range(emails){
			sqlStr += "(" + strconv.Itoa(userid) + ", ?), "
			vals = append(vals, dataSet[userid][i])
		}
		sqlStr = sqlStr[0:len(sqlStr)-2]
		stmt, err := db.Prepare(sqlStr)
		checkErr(err)
		_, err = stmt.Exec(vals...)
		checkErr(err)
	}
}

func (update *Update)Delete(dataSet map[int][]string){
	// Takes (int, string[])map of data and removes
	// listed items from database
	db := update.db
	for userid, emails := range dataSet{
		tableName := "unsub_" + strconv.Itoa(int(math.Mod(float64(userid),15.0)))
		sqlStr := "DELETE FROM " + tableName + " WHERE user_id = " + strconv.Itoa(userid) + " AND ("
		var vals []interface{}
		for i:= range(emails){
			sqlStr += "email = ? OR "
			vals = append(vals, dataSet[userid][i])
		}
		sqlStr = sqlStr[0:len(sqlStr)-4]
		sqlStr += ")"
		stmt, err := db.Prepare(sqlStr)
		checkErr(err)
		_, err = stmt.Exec(vals...)
		checkErr(err)
	}
}

func (update *Update)Clear(){
	// Delete all rows from all tables in the database
	// Be very careful when using this function! It can
	// take a while to repopulate the db
	db := update.db
	for i := 0; i < 15; i++ {
		_, err := db.Exec("TRUNCATE TABLE unsub_" + strconv.Itoa(i))
		checkErr(err)
	}
}

