package bloomManager

import (
	"Inf191BloomFilter/databaseAccessObj"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateBloomFilter(t *testing.T) {
	// testing initial population of bloom filter
	bf := New()
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	arrayOfEmails := []string{"test1@uci.edu", "test2@uci.edu", "test3@uci.edu"}
	databaseTestMap := make(map[int][]string)
	databaseTestMap[0] = arrayOfEmails
	dao.Insert(databaseTestMap)
	bf.UpdateBloomFilter()
	assert.Equal(t, true, bf.bloomFilter.TestString("0_test1@uci.edu"))
	assert.Equal(t, true, bf.bloomFilter.TestString("0_test2@uci.edu"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_test0@gmail.com"))
	assert.Equal(t, true, bf.bloomFilter.TestString("0_test3@uci.edu"))
	dao.Delete(databaseTestMap)
}

func TestUpdateBloomFilter2(t *testing.T) {
	// testing update bloom filter
	bf := New()
	bf.UpdateBloomFilter()
	assert.Equal(t, false, bf.bloomFilter.TestString("0_catlover@uci.edu"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_snowcone3@uci.edu"))
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	arrayOfEmails := []string{"catlover@uci.edu", "snowcone3@uci.edu"}
	databaseTestMap := make(map[int][]string)
	databaseTestMap[0] = arrayOfEmails
	dao.Insert(databaseTestMap)
	bf.UpdateBloomFilter()
	assert.Equal(t, true, bf.bloomFilter.TestString("0_catlover@uci.edu"))
	assert.Equal(t, true, bf.bloomFilter.TestString("0_snowcone3@uci.edu"))
	dao.Delete(databaseTestMap)
}

func TestRepopulateBloomFilter(t *testing.T) {
	// testing repopulating bloom filter
	bf := New()
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	arrayOfEmails := []string{"ilovepadthai@gmail.com", "eatmyshorts@yahoo.com"}
	databaseTestMap := make(map[int][]string)
	databaseTestMap[0] = arrayOfEmails
	// add two records to database before test
	dao.Insert(databaseTestMap)
	bf.UpdateBloomFilter()
	assert.Equal(t, true, bf.bloomFilter.TestString("0_ilovepadthai@gmail.com"))
	assert.Equal(t, true, bf.bloomFilter.TestString("0_eatmyshorts@yahoo.com"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_test2@gmail.com"))
	// delete two records from database to test resubscribe
	dao.Delete(databaseTestMap)
	bf.RepopulateBloomFilter()
	assert.Equal(t, false, bf.bloomFilter.TestString("0_ilovepadthai@gmail.com"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_eatmyshorts@yahoo.com"))
	assert.Equal(t, false, bf.bloomFilter.TestString("0_test2@gmail.com"))
}

func TestGetArrayOfUnsubscribedEmails(t *testing.T) {
	arrayOfEmails := []string{"0_ilovepadthai@gmail.com",
		"0_eatmyshorts@yahoo.com",
		"0_friedchicken@gmail.com",
		"0_juicebar@uci.edu",
		"0_ratatouille@hungry.com",
		"0_chocolatebar@yahoo.com"}
	bf := New()
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	arrayOfEmails2 := []string{"friedchicken@gmail.com", "chocolatebar@yahoo.com"}
	databaseTestMap := make(map[int][]string)
	databaseTestMap[0] = arrayOfEmails2
	dao.Insert(databaseTestMap)
	bf.UpdateBloomFilter()
	var emailToCheck []string
	emailToCheck = bf.GetArrayOfUnsubscribedEmails(arrayOfEmails)
	assert.Equal(t, 2, len(emailToCheck))
	assert.Equal(t, "0_friedchicken@gmail.com", emailToCheck[0])
	assert.Equal(t, "0_chocolatebar@yahoo.com", emailToCheck[1])
	dao.Delete(databaseTestMap)
}
