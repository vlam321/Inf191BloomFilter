package bloomManager

import (
	"Inf191BloomFilter/databaseAccessObj"
	"strconv"
	"strings"

	"github.com/willf/bloom"
)

const databaseSize = 15

// BloomFilter struct holds the pointer to the bloomFilter object
type BloomFilter struct {
	bloomFilter *bloom.BloomFilter
	bitArraySize uint
	numHashFunc uint
}

// New is called to instantiate a new BloomFilter object
func New(bitArraySize, numHashFunc uint) *BloomFilter {
	bloomFilter := bloom.New(bitArraySize, numHashFunc)
	return &BloomFilter{bloomFilter, bitArraySize, numHashFunc}
}

func (bf *BloomFilter) getStats(dbSize uint) (uint, uint, float64){
	return bf.bitArraySize, bf.numHashFunc, bf.bloomFilter.EstimateFalsePositiveRate(dbSize)
}

// UpdateBloomFilter will be called if unsubscribed emails are added to the database
// (unsubscribe emails), can be used for initially populating the bloom filter and
// updating the bloom filter
func (bf *BloomFilter) UpdateBloomFilter() {
	arrayOfUserIDEmail := getArrayOfUserIDEmail()
	for i := range arrayOfUserIDEmail {
		bf.bloomFilter.AddString(arrayOfUserIDEmail[i])
	}
}

// RepopulateBloomFilter will be called if unsubscribed emails are removed from the
// database (customers resubscribe to emails)
func (bf *BloomFilter) RepopulateBloomFilter() {
	newBloomFilter := bloom.New(bf.bitArraySize,bf.numHashFunc)
	arrayOfUserIDEmail := getArrayOfUserIDEmail()
	for i := range arrayOfUserIDEmail {
		newBloomFilter.AddString(arrayOfUserIDEmail[i])
	}
	bf.bloomFilter = newBloomFilter.Copy()
}

// getArrayOfUserIDEmail retrieves all records in the database shards and returns an array
// of strings in the form of userid_email
func getArrayOfUserIDEmail() []string {
	var arrayOfUserIDEmail []string
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	// loops through all tables in the database
	for j := 0; j < databaseSize; j++ {
		databaseResultMap := dao.SelectTable(j)
		for key, value := range databaseResultMap {
			for i := range value {
				arrayOfUserIDEmail = append(arrayOfUserIDEmail, strconv.Itoa(int(key))+"_"+value[i])
			}
		}
	}
	dao.CloseConnection()
	return arrayOfUserIDEmail
}

// GetArrayOfUnsubscribedEmails given a list of strings will return a list of those
// that exist in the bloom filter
func (bf *BloomFilter) GetArrayOfUnsubscribedEmails(arrayOfEmails []string) []string {
	var arrayOfUnsubscribedEmails []string
	var arrayOfUnsubscribedEmails2 []string
	mapOfPositives := make(map[int][]string)
	// filters true results into an array
	for i := range arrayOfEmails {
		if bf.bloomFilter.TestString(arrayOfEmails[i]) {
			arrayOfUnsubscribedEmails = append(arrayOfUnsubscribedEmails, arrayOfEmails[i])
		}
	}
	// convert to map to check all trues with database to filter out the false positives
	for i := range arrayOfUnsubscribedEmails {
		keyValueArray := strings.Split(arrayOfUnsubscribedEmails[i], "_")
		var key, _ = strconv.Atoi(keyValueArray[0])
		var value = keyValueArray[1]
		_, ok := mapOfPositives[key]
		if ok {
			mapOfPositives[key] = append(mapOfPositives[key], value)
		} else {
			var valueArray []string
			mapOfPositives[key] = append(valueArray, value)
		}
	}
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	databaseResultMap := dao.Select(mapOfPositives)
	// convert back to array
	for key, value := range databaseResultMap {
		for i := range value {
			arrayOfUnsubscribedEmails2 = append(arrayOfUnsubscribedEmails2, strconv.Itoa(int(key))+"_"+value[i])
		}
	}
	dao.CloseConnection()
	return arrayOfUnsubscribedEmails2
}
