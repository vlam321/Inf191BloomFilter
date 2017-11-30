package databaseAccessObj

import (
	"database/sql"
	"math"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Update struct {
	db *sql.DB
}

type Pair struct {
	// pair struct that holds user_id:[emails]
	id     int
	emails []string
}

type SqlStrVal struct {
	sqlStr string
	val    []interface{}
}

type GraphValue struct {
	graphType string
	x         float32
	y         float32
}

func modId(userid int) int {
	// mod user_id by 15
	return int(math.Mod(float64(userid), 15.0))
}

func New(dsn string) *Update {
	// input: dsn = 'username:password@/database'
	// create new connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return &Update{db}
}

func (update *Update) hasTable(databaseName, tableName string) bool {
	// For testing purposes
	// Checks if the specify tablename exist in the specified database
	db := update.db
	var check string
	err := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ? LIMIT 1", databaseName, tableName).Scan(&check)
	if err == sql.ErrNoRows {
		return false
	} else {
		checkErr(err)
	}
	return true
}

func (update *Update) dropTable(tableName string) {
	// for testing purposes
	// Delete the table in the db if exists
	db := update.db
	_, err := db.Exec("DROP TABLE IF EXISTS " + tableName)
	checkErr(err)
}

func (update *Update) CloseConnection() {
	// Update method, closes connection
	update.db.Close()
}

func (update *Update) MkTbl(tablename, schema string) {
	db := update.db
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tablename + schema + ";")
	checkErr(err)
}

func (update *Update) SelectRandSubset(tblNum, size int) map[int][]string {
	db := update.db
	result := make(map[int][]string)

	stmt, err := db.Prepare("SELECT user_id, email FROM unsub_" + strconv.Itoa(tblNum) + " ORDER BY RAND() LIMIT ?;")
	checkErr(err)

	rows, err := stmt.Query(strconv.Itoa(size))
	checkErr(err)

	defer rows.Close()
	var user_id int
	var email string
	for rows.Next() {
		err = rows.Scan(&user_id, &email)
		checkErr(err)
		result[user_id] = append(result[user_id], email)
	}
	return result
}

func (update *Update) Select(dataSet map[int][]string) map[int][]string {
	// Return items that exist both in input dataSet and database
	db := update.db
	result := make(map[int][]string)
	for userid, emails := range dataSet {
		tableName := "unsub_" + strconv.Itoa(modId(userid))
		sqlStr := "SELECT user_id, email FROM " + tableName + " WHERE user_id = " + strconv.Itoa(userid) + " AND ("
		var vals []interface{}
		for i := range emails {
			sqlStr += "email = ? OR "
			vals = append(vals, dataSet[userid][i])
		}
		sqlStr = sqlStr[0 : len(sqlStr)-4]
		sqlStr += ")"
		rows, err := db.Query(sqlStr, vals...)
		if err == sql.ErrNoRows {
			continue
		}
		checkErr(err)

		for rows.Next() {
			var user_id int
			var email string
			err = rows.Scan(&user_id, &email)
			checkErr(err)
			result[user_id] = append(result[user_id], email)
		}
	}
	return result
}

func (update *Update) Tselect(dataSet map[int][]string) map[int][]string {
	// Return items that exist both in input dataSet and database
	db := update.db
	result := make(map[int][]string)

	var sqlStrings []SqlStrVal

	for userid, emails := range dataSet {
		counter := 0
		tableName := "unsub_" + strconv.Itoa(modId(userid))
		sqlStr := "SELECT user_id, email FROM " + tableName + " WHERE user_id = " + strconv.Itoa(userid) + " AND ("
		var vals []interface{}
		for i := range emails {
			sqlStr += "email = ? OR "
			vals = append(vals, dataSet[userid][i])
			counter += 1
			if counter >= 64000 {
				sqlStrings = append(sqlStrings, SqlStrVal{sqlStr, vals[0:len(vals)]})
				sqlStr = "SELECT user_id, email FROM " + tableName + " WHERE user_id = " + strconv.Itoa(userid) + " AND ("
				vals = make([]interface{}, 0)
				counter = 0

			}
		}
		if len(vals) != 0 {
			sqlStrings = append(sqlStrings, SqlStrVal{sqlStr, vals[0:len(vals)]})
		}
		for i := range sqlStrings {
			rows, err := db.Query(sqlStrings[i].sqlStr[0:len(sqlStrings[i].sqlStr)-4]+")", sqlStrings[i].val[0:len(sqlStrings[i].val)]...)
			if err == sql.ErrNoRows {
				continue
			}
			checkErr(err)
			defer rows.Close()
			for rows.Next() {
				var user_id int
				var email string
				err = rows.Scan(&user_id, &email)
				checkErr(err)
				result[user_id] = append(result[user_id], email)
			}
		}
	}
	return result
}

func (update *Update) SelectByTimestamp(ts time.Time) map[int][]string {
	// Select all items from database where input time after item's timestamp
	db := update.db
	result := make(map[int][]string)
	for i := 0; i < 15; i++ {
		tableName := "unsub_" + strconv.Itoa(i)
		sqlStr := "SELECT user_id, email FROM " + tableName + " WHERE ts >= ?"
		rows, err := db.Query(sqlStr, ts.String())
		if err == sql.ErrNoRows {
			continue
		}
		checkErr(err)
		defer rows.Close()
		for rows.Next() {
			var user_id int
			var email string
			err = rows.Scan(&user_id, &email)
			checkErr(err)
			result[user_id] = append(result[user_id], email)
		}
	}
	return result
}

func (update *Update) SelectTable(tableNum int) map[int][]string {
	// Select all items from a single table
	db := update.db
	result := make(map[int][]string)
	tableName := "unsub_" + strconv.Itoa(tableNum)
	sqlStr := "SELECT user_id, email FROM " + tableName
	rows, err := db.Query(sqlStr)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var user_id int
		var email string
		err = rows.Scan(&user_id, &email)
		checkErr(err)
		result[user_id] = append(result[user_id], email)
	}
	return result
}

func (update *Update) Insert(dataSet map[int][]string) {
	// Takes (int, string[])map of data and inserts
	// listed items into database
	db := update.db
	shardMap := make(map[int][]Pair)
	for userid, emails := range dataSet {
		shardMap[modId(userid)] = append(shardMap[modId(userid)], Pair{userid, emails})
	}
	for tabNum, pairs := range shardMap {
		tableName := "unsub_" + strconv.Itoa(tabNum)
		sqlStr := "INSERT INTO " + tableName + "(user_id, email, ts) VALUES "
		var vals []interface{}
		counter := 0
		for p := range pairs {
			for e := range pairs[p].emails {
				sqlStr += "(?, ?, CURRENT_TIMESTAMP), "
				vals = append(vals, pairs[p].id, pairs[p].emails[e])
				counter += 1
			}
		}
		sqlStr = sqlStr[0 : len(sqlStr)-2]
		stmt, err := db.Prepare(sqlStr)
		checkErr(err)
		_, err = stmt.Exec(vals...)
		checkErr(err)
	}
}

func (update *Update) Tinsert(dataSet map[int][]string) {
	// Takes (int, string[])map of data and inserts
	// listed items into database
	db := update.db
	shardMap := make(map[int][]Pair)
	for userid, emails := range dataSet {
		shardMap[modId(userid)] = append(shardMap[modId(userid)], Pair{userid, emails})
	}

	var sqlStrings []SqlStrVal

	for tabNum, pairs := range shardMap {
		tableName := "unsub_" + strconv.Itoa(tabNum)
		sqlStr := "INSERT INTO " + tableName + "(user_id, email, ts) VALUES "
		var vals []interface{}
		counter := 0
		for p := range pairs {
			for e := range pairs[p].emails {
				sqlStr += "(?, ?, CURRENT_TIMESTAMP), "
				vals = append(vals, pairs[p].id, pairs[p].emails[e])
				counter += 1
				if counter >= 32000 {
					sqlStrings = append(sqlStrings, SqlStrVal{sqlStr, vals[0:len(vals)]})
					sqlStr = "INSERT INTO " + tableName + "(user_id, email, ts) VALUES "
					vals = make([]interface{}, 0)
					counter = 0
				}
			}
		}
		if len(vals) != 0 {
			sqlStrings = append(sqlStrings, SqlStrVal{sqlStr, vals[0:len(vals)]})
		}

	}
	for i := range sqlStrings {
		stmt, err := db.Prepare(sqlStrings[i].sqlStr[0 : len(sqlStrings[i].sqlStr)-2])
		checkErr(err)
		_, err = stmt.Exec(sqlStrings[i].val...)
		checkErr(err)
	}
}

func (update *Update) LogTestResult(resultType string, x, y float32) {
	db := update.db
	sqlStr := "INSERT INTO test_results (result_type, x_axis, y_axis) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(sqlStr)
	checkErr(err)
	_, err = stmt.Exec(resultType, x, y)
	checkErr(err)
}

func (update *Update) SelectTestResults() []GraphValue {
	db := update.db
	sqlStr := "SELECT result_type, x_axis, y_axis FROM test_results"
	stmt, err := db.Prepare(sqlStr)
	checkErr(err)
	rows, err := stmt.Query()
	checkErr(err)

	var resultType string
	var x float32
	var y float32
	result := make([]GraphValue, 0)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&resultType, &x, &y)
		checkErr(err)
		result = append(result, GraphValue{resultType, x, y})
	}
	return result[0:len(result)]

}

func (update *Update) Delete(dataSet map[int][]string) {
	// Takes (int, string[])map of data and removes
	// listed items from database
	db := update.db
	for userid, emails := range dataSet {
		tableName := "unsub_" + strconv.Itoa(modId(userid))
		sqlStr := "DELETE FROM " + tableName + " WHERE user_id = " + strconv.Itoa(userid) + " AND ("
		var vals []interface{}
		for i := range emails {
			sqlStr += "email = ? OR "
			vals = append(vals, dataSet[userid][i])
		}
		sqlStr = sqlStr[0 : len(sqlStr)-4]
		sqlStr += ")"
		stmt, err := db.Prepare(sqlStr)
		checkErr(err)
		_, err = stmt.Exec(vals...)
		checkErr(err)
	}
}

func (update *Update) Clear() {
	// Delete all rows from all tables in the database
	// Be very careful when using this function! It can
	// take a while to repopulate the db
	db := update.db
	for i := 0; i < 15; i++ {
		_, err := db.Exec("TRUNCATE TABLE unsub_" + strconv.Itoa(i))
		checkErr(err)
	}
}
