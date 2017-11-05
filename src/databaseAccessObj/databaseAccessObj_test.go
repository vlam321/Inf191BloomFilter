/* Primamry script for running unit tests for the databaseAccessObj.
For running running all the test cases, use the following command:
	- go test databaseAccessObj_test.go databaseAccessObj.go
Notice that the above command will not show any print statements or
detail about specific tests. If you want to see those, use the -v flag,
like so:
	- go test databaseAccessObj_test.go databaseAccessObj.go
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
	assert.Equal(t, update.HasTable("unsubscribed", "unsub_0"), true)
	assert.Equal(t, update.HasTable("unsubscribed", "unsub_1"), false)
	assert.Equal(t, update.HasTable("unsubscribed", "unsub_2"), false)
	update.CloseConnection()
}

