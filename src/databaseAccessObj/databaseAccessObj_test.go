/* Primamry script for running unit tests for the databaseAccessObj.
For running running all the test cases, use the following command:
	
	go test databaseAccessObj_test.go databaseAccessObj.go

Notice that the above command will not show any print statements or
detail about specific tests. If you want to see those, use the -v flag,
like so:
	
	go test -v databaseAccessObj_test.go databaseAccessObj.go
*/

package databaseAccessObj

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const dsn = "bloom:test@/unsubscribed"

func TestConnection(t *testing.T){
	update := New(dsn)
	update.CloseConnection()
}

func TestHasTable(t *testing.T){
	update := New(dsn)
	assert.Equal(t, update.hasTable("unsubscribed", "unsub_0"), true)
	assert.Equal(t, update.hasTable("unsubscribed", "nothere"), false)
	assert.Equal(t, update.hasTable("unsubscribed", "nothereeither"), false)
	assert.Equal(t, update.hasTable("unsubscribed", "nope"), false)
	update.CloseConnection()
}


func TestEnsureTable(t *testing.T){
	update := New(dsn)
	// Test ensure table on a table that does exist to make sure
	// there's no bugs
	assert.Equal(t, update.hasTable("unsubscribed", "unsub_0"), true)

	// Delete these tables if they exist in the db
	update.dropTable("unsub_1")
	update.dropTable("unsub_2")

	// Make sure that these table doesn't exist in the databse
	assert.Equal(t, update.hasTable("unsubscribed", "unsub_1"), false)
	assert.Equal(t, update.hasTable("unsubscribed", "unsub_2"), false)

	// Run ensure table using the same inputs to create the tables
	update.EnsureTable("unsub_1")
	update.EnsureTable("unsub_2")

	// double check to make sure that the tables are created
	assert.Equal(t, update.hasTable("unsubscribed", "unsub_1"), true)
	assert.Equal(t, update.hasTable("unsubscribed", "unsub_2"), true)

	update.CloseConnection()
}
