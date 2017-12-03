package bloomManager

import (
	"testing"

	"Inf191BloomFilter/databaseAccessObj"

	"github.com/stretchr/testify/assert"
)

const bitArraySize = 100000
const numHash = 5

/*
func TestUpdateBloomFilter(t *testing.T) {
	// testing initial population of bloom filter
	bf := New(bitArraySize,numHash)
	db := databaseAccessObj.New()
	databaseTestMap := make(map[int][]string)
	databaseTestMap[0] := []string{"test1@uci.edu", "test2@uci.edu", "test3@uci.edu"}

	// clean out database of test emails
	db.Delete(databaseTestMap)
	db.Insert(databaseTestMap)
	bf.UpdateBloomFilter()
	assert.Equal(t, true, bf.bloomFilter.TestString("0_test1@uci.edu"))
	assert.Equal(t, true, bf.bloomFilter.TestString("0_test2@uci.edu"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_test0@gmail.com"))
	assert.Equal(t, true, bf.bloomFilter.TestString("0_test3@uci.edu"))
	db.Delete(databaseTestMap)
	db.CloseConnection()
}
*/

/*
func TestUpdateBloomFilter2(t *testing.T) {
	// testing update bloom filter
	bf := New()
	db := databaseAccessObj.New("bloom:test@/unsubscribed")
	arrayOfEmails := []string{"catlover@uci.edu", "snowcone3@uci.edu"}
	databaseTestMap := make(map[int][]string)
	databaseTestMap[0] = arrayOfEmails
	// clean out database of test emails
	db.Delete(databaseTestMap)
	bf.UpdateBloomFilter()
	assert.Equal(t, false, bf.bloomFilter.TestString("0_catlover@uci.edu"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_snowcone3@uci.edu"))
	db.Insert(databaseTestMap)
	bf.UpdateBloomFilter()
	assert.Equal(t, true, bf.bloomFilter.TestString("0_catlover@uci.edu"))
	assert.Equal(t, true, bf.bloomFilter.TestString("0_snowcone3@uci.edu"))
	db.Delete(databaseTestMap)
	db.CloseConnection()
}
*/

func TestRepopulateBloomFilter(t *testing.T) {
	// testing repopulating bloom filter
	bf := New(bitArraySize, numHash)
	db := databaseAccessObj.New()
	defer db.CloseConnection()
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
	arrayOfEmails2 := []string{"friedchicken@gmail.com", "chocolatebar@yahoo.com"}
	databaseTestMap := make(map[int][]string)
	databaseTestMap[0] = arrayOfEmails2
	// clean out database of test emails
	db.Delete(databaseTestMap)
	db.Insert(databaseTestMap)
	bf.RepopulateBloomFilter()
	emailToCheck := bf.GetArrayOfUnsubscribedEmails(testMap)
	assert.Equal(t, 2, len(emailToCheck[0]))
	assert.Equal(t, "friedchicken@gmail.com", emailToCheck[0][1])
	assert.Equal(t, "chocolatebar@yahoo.com", emailToCheck[0][0])
	db.Delete(databaseTestMap)
}
