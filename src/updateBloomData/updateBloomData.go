package updateBloomData

import(
	//"../bloomDataGenerator"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

func outErr(err error){
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

func (update *Update)CloseConnection(){
	// Update method, closes connection
	update.db.Close()
	fmt.Println("Connection Closed\n")
}

func (update *Update)InsertDataSet(dataSet map[int][]string){
	// Takes a (int, string[])map of data and insert them
	// into the specified db
	db := update.db
	stmt, err := db.Prepare(`INSERT INTO unsub_0 (user_id, email) VALUES (?,?)`)
	outErr(err)

	_, err = db.Exec("BEGIN")
	outErr(err)

	for userid, emails := range dataSet{
		for i := range(emails){
			_, err := stmt.Exec(userid, emails[i])
			outErr(err)
		}
	}

	_, err = db.Exec("COMMIT")
	outErr(err)
}

func (update *Update)Clear() {
	// Delete all rows from a table in the database
	db := update.db
	_, err := db.Exec("TRUNCATE TABLE unsub_0")
	outErr(err)
}
