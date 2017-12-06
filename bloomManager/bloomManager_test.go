package bloomManager

import (
	"testing"

	"Inf191BloomFilter/databaseAccessObj"

	"github.com/stretchr/testify/assert"
)

const bitArraySize = 100000
const numHash = 10

func TestRepopulateBloomFilter(t *testing.T) {
	// testing repopulating bloom filter
	bf := New(bitArraySize, numHash)
	db := databaseAccessObj.New()
	defer db.CloseConnection()
	db.Clear()
	arrayOfEmails := []string{"ilovepadthai@gmail.com", "eatmyshorts@yahoo.com"}
	databaseTestMap := make(map[int][]string)
	databaseTestMap[0] = arrayOfEmails
	// clean out database of test emails
	db.Delete(databaseTestMap)
	// add two records to database before test
	db.Insert(databaseTestMap)
	bf.RepopulateBloomFilter()
	assert.Equal(t, true, bf.bloomFilter.TestString("0_ilovepadthai@gmail.com"))
	assert.Equal(t, true, bf.bloomFilter.TestString("0_eatmyshorts@yahoo.com"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_test2@gmail.com"))
	// delete two records from database to test resubscribe
	db.Delete(databaseTestMap)
	bf.RepopulateBloomFilter()
	assert.Equal(t, false, bf.bloomFilter.TestString("0_ilovepadthai@gmail.com"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_eatmyshorts@yahoo.com"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_test2@gmail.com"))
}

func TestGetArrayOfUnsubscribedEmails(t *testing.T) {
	arrayOfEmails := []string{"ilovepadthai@gmail.com",
		"eatmyshorts@yahoo.com",
		"friedchicken@gmail.com",
		"juicebar@uci.edu",
		"ratatouille@hungry.com",
		"chocolatebar@yahoo.com"}
	testMap := map[int][]string{0:arrayOfEmails}
	bf := New(bitArraySize, numHash)
	db := databaseAccessObj.New()
	defer db.CloseConnection()
	db.Clear()
	arrayOfEmails2 := []string{"friedchicken@gmail.com", "chocolatebar@yahoo.com"}
	databaseTestMap := make(map[int][]string)
	databaseTestMap[0] = arrayOfEmails2
	db.Insert(databaseTestMap)
	bf.RepopulateBloomFilter()
	emailToCheck := bf.GetArrayOfUnsubscribedEmails(testMap)
	assert.Equal(t, 2, len(emailToCheck[0]))
	assert.Equal(t, arrayOfEmails2[0], emailToCheck[0][0])
	assert.Equal(t, arrayOfEmails2[1], emailToCheck[0][1])
	db.Delete(databaseTestMap)
}
