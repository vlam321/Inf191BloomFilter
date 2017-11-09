package updateBloomData

import (
	"Inf191BloomFilter/src/databaseAccessObj"
	"strconv"

	"github.com/willf/bloom"
)

type BloomFilter struct {
	bloomFilter           bloom.BloomFilter
	backupCopyBloomFilter bloom.BloomFilter
}

func CopyAndRepopulateBloomFilter(bf bloom.BloomFilter) bloom.BloomFilter {
	// makes a copy of the old bloom filter and repopulates the bloom filter with new data
	// used when the database is modified
	backupCopyBloomFilter := bf.Copy()
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	databaseResultMap := dao.SelectAll()
	//access dao and selectall, return map
	bf.ClearAll()
	for key, value := range databaseResultMap {
		for i := 0; i <= len(value); i++ {
			bf.AddString(strconv.Itoa(int(key)) + "_" + value[i])
		}
	}
	bloomFilter := bf.Copy()
	return bf
}
