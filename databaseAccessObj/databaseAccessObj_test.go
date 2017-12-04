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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	db := New()
	db.CloseConnection()
}

func TestHasTable(t *testing.T) {
	db := New()
	db.Clear()
	assert.True(t, db.hasTable("unsubscribed", "unsub_0"))
	assert.False(t, db.hasTable("unsubscribed", "nothere"))
	assert.False(t, db.hasTable("unsubscribed", "nothereeither"))
	assert.False(t, db.hasTable("unsubscribed", "nope"))
	db.CloseConnection()
}

func TestInsertAndSelect(t *testing.T) {
	db := New()
	db.Clear()
	testData := make(map[int][]string)
	testData[0] = []string{"a", "b", "c"}
	testData[16] = []string{"d", "e", "f", "g"}
	db.Insert(testData)
	result := db.Select(testData)
	assert.Equal(t, testData, result)
	db.CloseConnection()
}

func TestSelectTable(t *testing.T) {
	db := New()
	db.Clear()
	testData := make(map[int][]string)
	testData[3] = []string{"h", "i", "j"}
	testData[18] = []string{"k", "l", "m", "n", "o"}
	testData[23] = []string{"p", "q", "r", "s"}
	testDataShard3 := make(map[int][]string)
	testDataShard3[3] = testData[3]
	testDataShard3[18] = testData[18]
	testDataShard8 := make(map[int][]string)
	testDataShard8[23] = testData[23]
	db.Insert(testData)
	result := db.SelectTable(3)
	assert.Equal(t, testDataShard3, result)
	result = db.SelectTable(8)
	assert.Equal(t, testDataShard8, result)
	db.CloseConnection()
}

func TestSelectByTimestamp(t *testing.T) {
	db := New()
	db.Clear()
	testData := make(map[int][]string)
	testData[6] = []string{"t", "u", "v", "w"}
	testData2 := make(map[int][]string)
	testData2[6] = []string{"x", "y", "z"}
	db.Insert(testData)
	time.Sleep(time.Second)
	ts := time.Now()
	time.Sleep(time.Second)
	db.Insert(testData2)
	result := db.SelectByTimestamp(ts)
	assert.Equal(t, testData2, result)
	db.CloseConnection()
}
