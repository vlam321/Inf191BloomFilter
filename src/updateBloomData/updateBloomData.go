package updateBloomData

import (
	"Inf191BloomFilter/src/databaseAccessObj"
	"strconv"

	"github.com/willf/bloom"
)

const bitArraySize = 10000
const numberOfHashFunction = 5

type BloomFilter struct {
	bloomFilter *bloom.BloomFilter
}

func New() *BloomFilter {
	bloomFilter := bloom.New(bitArraySize, numberOfHashFunction)
	return &BloomFilter{bloomFilter}
}

func (bf *BloomFilter) UpdateBloomFilter() {
	// used when more unsubscribed emails have been added to the database
}

func (bf *BloomFilter) RepopulateBloomFilter() {
	// used when unsubscribed emails are removed from the database - resubscribed emails example
	newBloomFilter := bloom.New(bitArraySize, numberOfHashFunction)
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	databaseResultMap := dao.SelectAll()
	for key, value := range databaseResultMap {
		for i := range value {
			newBloomFilter.AddString(strconv.Itoa(int(key)) + "_" + value[i])
		}
	}
	bf.bloomFilter = newBloomFilter.Copy()
}
