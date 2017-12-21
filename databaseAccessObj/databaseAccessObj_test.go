// Primary script for running unit tests for the databaseAccessObj.

package databaseAccessObj

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	db := New()
	defer db.CloseConnection()
}

func TestHasTable(t *testing.T) {
	db := New()
	defer db.CloseConnection()
	db.Clear()
	assert.True(t, db.hasTable("unsubscribed", "unsub_0"))
	assert.False(t, db.hasTable("unsubscribed", "nothere"))
	assert.False(t, db.hasTable("unsubscribed", "nothereeither"))
	assert.False(t, db.hasTable("unsubscribed", "nope"))
}

func TestInsertAndSelect(t *testing.T) {
	db := New()
	defer db.CloseConnection()
	db.Clear()
	testData := make(map[int][]string)
	testData[0] = []string{"a", "b", "c"}
	testData[16] = []string{"d", "e", "f", "g"}
	db.Insert(testData)
	result := db.Select(testData)
	assert.Equal(t, testData, result)
}

func TestSelectRandSubset(t *testing.T) {
	db := New()
	defer db.CloseConnection()
	db.Clear()
	testData := make(map[int][]string)
	testData[0] = [] string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	db.Insert(testData)
	result := db.SelectRandSubset(0, 4)
	assert.Equal(t, 4, len(result[0]))
}

func TestSelectTable(t *testing.T) {
	db := New()
	defer db.CloseConnection()
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
}

func TestSelectByTimestamp(t *testing.T) {
	db := New()
	defer db.CloseConnection()
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
}

