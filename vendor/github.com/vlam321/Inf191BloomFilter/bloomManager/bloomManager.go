package bloomManager

import (
	"log"
	"strconv"

	"github.com/vlam321/Inf191BloomFilter/databaseAccessObj"

	"github.com/willf/bloom"
)

const dbShards = 15

// BloomFilter struct holds the pointer to the bloomFilter object
type BloomFilter struct {
	bloomFilter *bloom.BloomFilter
}

// New is called to instantiate a new BloomFilter object
func New(numEmail uint, fpProb float64) *BloomFilter {
	bloomFilter := bloom.NewWithEstimates(numEmail, fpProb)
	log.Printf("BLOOM STATS: %d HASH FUNCTIONS | BIT ARRAY LEN OF %d", bloomFilter.K(), bloomFilter.Cap())
	return &BloomFilter{bloomFilter}
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
func (bf *BloomFilter) RepopulateBloomFilter(tableNum int) {
	db := databaseAccessObj.New()
	defer db.CloseConnection()
	numEmail := uint(db.GetTableSize(tableNum))
	newBloomFilter := bloom.NewWithEstimates(numEmail, float64(0.001))

	data := db.SelectTable(tableNum)
	for userid, emails := range data {
		u := strconv.Itoa(userid)
		for e := range emails {
			newBloomFilter.AddString(u + "_" + emails[e])
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
