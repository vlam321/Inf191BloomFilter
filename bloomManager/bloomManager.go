package bloomManager

import (
	"strconv"

	"Inf191BloomFilter/databaseAccessObj"

	"github.com/willf/bloom"
)

const dbShards = 15

// BloomFilter struct holds the pointer to the bloomFilter object
type BloomFilter struct {
	bloomFilter  *bloom.BloomFilter
	bitArraySize uint
	numHashFunc  uint
}

// New is called to instantiate a new BloomFilter object
func New(bitArraySize, numHashFunc uint) *BloomFilter {
	bloomFilter := bloom.New(bitArraySize, numHashFunc)
	return &BloomFilter{bloomFilter, bitArraySize, numHashFunc}
}

// GetStats returns false positive rate of bloom filter based on input size
func (bf *BloomFilter) GetStats(dbSize uint) float64 {
	return bf.bloomFilter.EstimateFalsePositiveRate(dbSize)
}

// not in use yet
// UpdateBloomFilter will be called if emails are added to the database
// (unsubscribe emails), can be used for initially populating the bloom filter and
// updating the bloom filter
/*
func (bf *BloomFilter) UpdateBloomFilter(ts time.Time) {
	db := databaseAccessObj.New()
	defer db.CloseConnection()

	for i := 0; i<dbShards; i++{
		data := db.SelectTable(i)
		for userid, emails := range(data){
			u := strconv.Itoa(userid)
			for e := range emails {
				bf.bloomFilter.AddString(u+"_"+emails[e])
			}
		}
	}
}
*/

// RepopulateBloomFilter will be called if emails are removed from the
// database (customers resubscribe)
// also used to initially populate bloom filter
func (bf *BloomFilter) RepopulateBloomFilter() {
	db := databaseAccessObj.New()
	defer db.CloseConnection()
	newBloomFilter := bloom.New(bf.bitArraySize, bf.numHashFunc)

	for i := 0; i < dbShards; i++ {
		data := db.SelectTable(i)
		for userid, emails := range data {
			u := strconv.Itoa(userid)
			for e := range emails {
				newBloomFilter.AddString(u + "_" + emails[e])
			}
		}
	}
	bf.bloomFilter = newBloomFilter.Copy()
}

// filter given a map[int][]string returns items that return true from bf
func (bf *BloomFilter) filter(dataSet map[int][]string) map[int][]string {
	result := make(map[int][]string)
	for userid, emails := range dataSet {
		u := strconv.Itoa(userid)
		for e := range emails {
			if bf.bloomFilter.TestString(u + "_" + emails[e]) {
				result[userid] = append(result[userid], emails[e])
			}
		}
	}
	return result
}

// GetArrayOfUnsubscribedEmails given a map of user_id:[emails] will return a map[user_id]:[emails]
// of those that exist in the db
func (bf *BloomFilter) GetArrayOfUnsubscribedEmails(dataSet map[int][]string) map[int][]string {
	// filters true results into an map[int][]string
	filtered := bf.filter(dataSet)
	db := databaseAccessObj.New()
	defer db.CloseConnection()
	result := db.Select(filtered)
	return result
}
