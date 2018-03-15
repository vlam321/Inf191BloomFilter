package databaseAccessObj

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// dbShards number of shards in database
const dbShards int = 15

// Update struct that holds db object
type Conn struct {
	db *sql.DB
}

// Pair struct that holds user_id:[emails]
type Pair struct {
	id     int
	emails []string
}

// SqlStrVal struct used to build SQL queries
type SqlStrVal struct {
	sqlStr string
	val    []interface{}
}

// Metrics struct holds metrics
type Metrics struct {
	result_type string
	X           float64
	Y           float64
}

// modId mod userid by number of database shards
func modId(userid int) int {
	return int(math.Mod(float64(userid), float64(dbShards)))
}

// New construct Conn object
func New() *Conn {
	viper.SetConfigName("sqlConn")
	viper.AddConfigPath("settings")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Database access object: %v\n", err.Error())
	}

	cfg := mysql.Config{
		//Addr: "mysql:3306",
		Addr:   viper.GetString("Addr"),
		User:   viper.GetString("User"),
		Passwd: viper.GetString("Passwd"),
		Net:    viper.GetString("Net"),
		DBName: viper.GetString("DBName"),
	}

	// log.Println("USING DSN = ", cfg.FormatDSN())
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		return nil
	}

	return &Conn{db}
}

// hasTable checks if the table tableName exists in db databaseName
func (conn *Conn) hasTable(databaseName, tableName string) bool {
	db := conn.db
	var check string
	err := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ? LIMIT 1", databaseName, tableName).Scan(&check)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Printf("Error : %v\n", err)
			return false
		}
	}

	return true
}

// dropTable removes table tableName from the db if it exists
func (conn *Conn) dropTable(tableName string) {
	db := conn.db
	_, err := db.Exec("DROP TABLE IF EXISTS " + tableName)

	if err != nil {
		log.Printf("Error dropping table: %v\n", err)
		return
	}
}

// CloseConnection closes connection to db
func (conn *Conn) CloseConnection() {
	conn.db.Close()
}

// MakeTable creates table tableName with attributes schema
func (conn *Conn) MakeTable(tablename, schema string) {
	db := conn.db
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tablename + schema + ";")

	if err != nil {
		log.Printf("Error creating table: %v\n", err)
		return
	}
}

// SelectRandSubset selects a random subset of data from shard tblNum in db
func (conn *Conn) SelectRandSubset(tblNum, size int) map[int][]string {
	db := conn.db
	result := make(map[int][]string)

	sqlStr := "SELECT user_id, email FROM unsub_" + strconv.Itoa(tblNum) + " ORDER BY RAND() LIMIT ?"

	rows, err := db.Query(sqlStr, strconv.Itoa(size))
	if err != nil {
		log.Printf("Error query: %v\n", err)
		return nil
	}

	defer rows.Close()
	var user_id int
	var email string
	for rows.Next() {
		err = rows.Scan(&user_id, &email)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
		}
		result[user_id] = append(result[user_id], email)
	}
	return result
}

// Select returns items in database matching input dataSet
func (conn *Conn) Select(dataSet map[int][]string) map[int][]string {
	db := conn.db
	result := make(map[int][]string)

	for userid, emails := range dataSet {
		tableName := "unsub_" + strconv.Itoa(modId(userid))
		//sqlStr := "SELECT email FROM " + tableName + " WHERE user_id = ? and email = ?"
		sqlStr := fmt.Sprintf("SELECT email FROM %s WHERE user_id = ? and email IN (%s)",
			tableName,
			fmt.Sprintf("?"+strings.Repeat(",?", len(emails)-1)))
		args := make([]interface{}, len(emails)+1)
		args[0] = userid
		for i, email := range emails {
			args[i+1] = email
		}
		rows, err := db.Query(sqlStr, args...)
		if err != nil {
			log.Printf("Error querying db: %v\n", err)
		}
		var email string
		for rows.Next() {
			err = rows.Scan(&email)
			if err != nil {
				log.Printf("Error scanning row: %v\n", err)
			}
			result[userid] = append(result[userid], email)
		}
		defer rows.Close()
		/*
			if err != nil {
				log.Printf("Error preparing statement", err)
				return nil
			}
			for e := range emails {
				var user_id int
				var email string

				err = stmt.QueryRow(userid, emails[e]).Scan(&user_id, &email)
				if err != nil {
					if err == sql.ErrNoRows {
						continue
					} else {
						log.Printf("Error querying row", err)
						return nil
					}
				}
				result[user_id] = append(result[user_id], email)
			}
		*/
	}
	return result
}

// SelectByTimestamp returns items in database added at time ts
func (conn *Conn) SelectByTimestamp(ts string, tableNum int) map[int][]string {
	db := conn.db
	result := make(map[int][]string)

	tableName := "unsub_" + strconv.Itoa(tableNum)
	sqlStr := "SELECT user_id, email FROM " + tableName + " WHERE ts = ?"
	rows, err := db.Query(sqlStr, ts)
	if err != nil {
		log.Printf("Error query: %v\n", err)
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		var user_id int
		var email string
		err = rows.Scan(&user_id, &email)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return nil
		}
		result[user_id] = append(result[user_id], email)
	}
	return result
}

// SelectTable selects all data in table unsub_tableNum
// used for repopulating bloom filter
func (conn *Conn) SelectTable(tableNum int) map[int][]string {
	db := conn.db
	result := make(map[int][]string)
	tableName := "unsub_" + strconv.Itoa(tableNum)
	sqlStr := "SELECT user_id, email FROM " + tableName

	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Printf("Error query: %v\n", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var user_id int
		var email string
		err = rows.Scan(&user_id, &email)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
		}
		result[user_id] = append(result[user_id], email)
	}
	return result
}

func (conn *Conn) InsertToTable(tableNum int, dataSet []string) {
	db := conn.db
	tableName := "unsub_" + strconv.Itoa(tableNum)
	sqlStr := "INSERT INTO " + tableName + "(user_id, email, ts) VALUES "
	var vals []interface{}
	counter := 0

	var sqlStrings []SqlStrVal

	for i := range dataSet {
		sqlStr += "(?, ?, CURRENT_TIMESTAMP), "
		vals = append(vals, tableNum, dataSet[i])
		counter += 1
		if counter >= 32000 {
			sqlStrings = append(sqlStrings, SqlStrVal{sqlStr, vals[0:len(vals)]})
			sqlStr = "INSERT INTO " + tableName + "(user_id, email, ts) VALUES "
			vals = make([]interface{}, 0)
			counter = 0
		}
	}
	if len(vals) != 0 {
		sqlStrings = append(sqlStrings, SqlStrVal{sqlStr, vals[0:len(vals)]})
	}

	for i := range sqlStrings {
		stmt, err := db.Prepare(sqlStrings[i].sqlStr[0 : len(sqlStrings[i].sqlStr)-2])
		if err != nil {
			log.Printf("Error preparing statement: %v\n", err)
			return
		}

		_, err = stmt.Exec(sqlStrings[i].val...)
		if err != nil {
			log.Printf("Error executing statement: %v\n", err)
			return
		}
	}
}

// Insert inserts dataSet into db
func (conn *Conn) Insert(dataSet map[int][]string) {
	db := conn.db
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
		if err != nil {
			log.Printf("Error preparing statement: %v\n", err)
			return
		}
		_, err = stmt.Exec(sqlStrings[i].val...)
		if err != nil {
			log.Printf("Error executing statement: %v\n", err)
			return
		}
	}
}

// LogTestResult logs test result into db
func (conn *Conn) LogTestResult(resultType string, x, y float64) {
	db := conn.db
	sqlStr := "INSERT INTO test_results (result_type, x_axis, y_axis) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		log.Printf("Error preparing statement: %v\n", err)
		return
	}
	_, err = stmt.Exec(resultType, x, y)
	if err != nil {
		log.Printf("Error executing statemetn: %v\n", err)
		return
	}
}

// SelectTestResults selects test results from db
func (conn *Conn) SelectTestResults() []Metrics {
	db := conn.db
	sqlStr := "SELECT result_type, x_axis, y_axis FROM test_results"
	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Printf("Error query: %v\n", err)
	}

	var resultType string
	var x float64
	var y float64
	result := make([]Metrics, 0)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&resultType, &x, &y)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return nil
		}
		result = append(result, Metrics{resultType, x, y})
	}
	return result
}

func (conn *Conn) Delete(tableNum, rows int) {
	db := conn.db
	tableName := "unsub_" + strconv.Itoa(modId(tableNum))
	_, err := db.Exec("DELETE FROM " + tableName + " WHERE user_id=" + strconv.Itoa(tableNum) + " LIMIT " + strconv.Itoa(rows) + ";")
	if err != nil {
		log.Printf("Error clearing tables: %v\n", err)
		return
	}
}

/*
// Delete removes all items in db matching to dataSet
func (conn *Conn) Delete(dataSet map[int][]string) {
	db := conn.db
	for userid, emails := range dataSet {
		tableName := "unsub_" + strconv.Itoa(modId(userid))
		sqlStr := "DELETE FROM " + tableName + " WHERE user_id = " + strconv.Itoa(userid) + " AND ("
		var vals []interface{}

		for i := range emails {
			sqlStr += "email = ? OR "
			vals = append(vals, dataSet[userid][i])
		}
		sqlStr = sqlStr[0:len(sqlStr)-4] + ")"

		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			log.Printf("Error preparing statement: %v\n", err)
			return
		}
		_, err = stmt.Exec(vals...)
		if err != nil {
			log.Printf("Error executing statement %v\n", err)
			return
		}
	}
}
*/

// Get the size or number of rows of the table
func (conn *Conn) GetTableSize(tableNum int) int {
	db := conn.db
	sqlStr := "SELECT COUNT(*) FROM unsub_" + strconv.Itoa(tableNum) + ";"
	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Printf("Error: Unable to query count. %v\n", err.Error())
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Printf("Error: Unable to scan row counts %v\n", err.Error())
		}
	}
	return count
}

// Clear removes ALL rows from ALL tables in db
func (conn *Conn) Clear() {
	db := conn.db
	for i := 0; i < 15; i++ {
		_, err := db.Exec("TRUNCATE TABLE unsub_" + strconv.Itoa(i))
		if err != nil {
			log.Printf("Error clearing tables: %v\n", err)
			return
		}
	}
}

// ClearTestResults removes all test results from db
func (conn *Conn) ClearTestResults() {
	db := conn.db
	_, err := db.Exec("TRUNCATE TABLE test_results")
	if err != nil {
		log.Printf("Error clearing test results: %v\n", err)
		return
	}
}

// GetCountByTimestamp gets all unique timestamps and their cout from the db
func (conn *Conn) GetCountByTimestamp(tableNum int) map[string]int {
	db := conn.db
	sqlStr := "SELECT ts, count(*) from unsub_" + strconv.Itoa(tableNum) + " group by ts;"
	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Printf("Error: Unable to query timestamp by count. %v\n", err.Error())
	}
	defer rows.Close()

	result := make(map[string]int)

	for rows.Next() {
		var timestamp []uint8
		var count int
		err = rows.Scan(&timestamp, &count)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
		}
		result[string(timestamp)] = count
	}
	return result
}
