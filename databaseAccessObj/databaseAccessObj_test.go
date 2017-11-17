/* Primary script for running unit tests for the databaseAccessObj.
For running all the test cases, use the following command:
	
	go test databaseAccessObj_test.go databaseAccessObj.go

The above command will not show any print statements or
detail about specific tests. If you want to see those, use the -v flag,
like so:
	
	go test -v databaseAccessObj_test.go databaseAccessObj.go
*/

package databaseAccessObj

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"Inf191BloomFilter/bloomDataGenerator"
)

const dsn = "bloom:test@/unsubscribed"

func TestConnection(t *testing.T){
	update := New(dsn)
	update.CloseConnection()
}

func TestHasTable(t *testing.T){
	update := New(dsn)
	assert.True(t, update.hasTable("unsubscribed", "unsub_0"))
	assert.False(t, update.hasTable("unsubscribed", "nothere"))
	assert.False(t, update.hasTable("unsubscribed", "nothereeither"))
	assert.False(t, update.hasTable("unsubscribed", "nope"))
	update.CloseConnection()
}

func TestInsertAndSelect(t *testing.T){
	update := New(dsn)
	data := bloomDataGenerator.GenData(1, 10, 20)

	// Make sure the table is empty
	update.Clear()

	// Add data to db
	update.InsertDataSet(data)

	retrieved := update.SelectAll()
	expected_emails := make(map[string]bool)

	for _, emails := range data{
		for i := range emails{
			expected_emails[emails[i]] = true
		}
	}

	// Test if data generated is the same as ones
	// retrieved from db
	for _, emails := range retrieved {
		for i := range emails{
			if ! expected_emails[emails[i]]{
				t.Error()
			}
		}
	}
}
