package updateBloomData

import (
	"Inf191BloomFilter/databaseAccessObj"
	"strconv"

	"github.com/willf/bloom"
)

const bitArraySize = 10000
const numberOfHashFunction = 5

// BloomFilter struct holds the pointer to the bloomFilter object
type BloomFilter struct {
	bloomFilter *bloom.BloomFilter
}

// New is called to instantiate a new BloomFilter object
func New() *BloomFilter {
	bloomFilter := bloom.New(bitArraySize, numberOfHashFunction)
	return &BloomFilter{bloomFilter}
}

// UpdateBloomFilter is used when more unsubscribed emails have been added to the database
func (bf *BloomFilter) UpdateBloomFilter() {
	var arrayOfUserIDEmail = getArrayOfUserIDEmail()
	for i := range arrayOfUserIDEmail {
		bf.bloomFilter.AddString(arrayOfUserIDEmail[i])
	}
}

// RepopulateBloomFilter will be called if unsubscribed emails are removed from the
// database (customers resubscribe to emails)
func (bf *BloomFilter) RepopulateBloomFilter() {
	newBloomFilter := bloom.New(bitArraySize, numberOfHashFunction)
	var arrayOfUserIDEmail []string
	arrayOfUserIDEmail = getArrayOfUserIDEmail()
	for i := range arrayOfUserIDEmail {
		newBloomFilter.AddString(arrayOfUserIDEmail[i])
	}
	bf.bloomFilter = newBloomFilter.Copy()
}

// getArrayOfUserIDEmail retrieves all records in the database and returns an array
// of strings in the form of userid_email
func getArrayOfUserIDEmail() []string {
	var arrayOfUserIDEmail []string
	dao := databaseAccessObj.New("bloom:test@/unsubscribed")
	databaseResultMap := dao.SelectAll()
	for key, value := range databaseResultMap {
		for i := range value {
			arrayOfUserIDEmail = append(arrayOfUserIDEmail, strconv.Itoa(int(key))+"_"+value[i])
		}
	}
	dao.CloseConnection()
	return arrayOfUserIDEmail
}

//GetArrayOfUnsubscribedEmails given a list of strings will return a list of those
//that exist in the bloom filter
func (bf *BloomFilter) GetArrayOfUnsubscribedEmails(arrayOfEmails []string) []string {
	var arrayOfUnsubscribedEmails []string
	for i := range arrayOfEmails {
		if bf.bloomFilter.TestString(arrayOfEmails[i]) {
			arrayOfUnsubscribedEmails = append(arrayOfUnsubscribedEmails, arrayOfEmails[i])
		}
	}
	return arrayOfUnsubscribedEmails
}
