package bloomDataGenerator

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDataSize(t *testing.T){
	// Test the data size produced by GenData matche specification
	expected_num_users := 10
	expected_min_emails := 1000
	expected_max_emails := 10000

	data := GenData(expected_num_users, expected_min_emails, expected_max_emails)

	assert.Equal(t, len(data), expected_num_users)
	for _, emails := range data{
		assert.Equal(t, len(emails) >= expected_min_emails, true)
		assert.Equal(t, len(emails) <= expected_max_emails, true)
	}
}

func TestEmailAddressCharLength(t * testing.T){
	// Test that the averge email address is between 21-23
	// chars in length
	var charSum int
	var numEmails int
	dataset1 := GenData(1, 1000, 2000)
	dataset2 := GenData(1, 100000, 200000)
	for _, emails := range dataset1{
		for i := range emails{
			charSum += len(emails[i])
		}
		numEmails += len(emails)
	}
	assert.Equal(t, charSum/numEmails >= 21, true)
	assert.Equal(t, charSum/numEmails <= 23, true)

	charSum = 0
	numEmails = 0
	for _, emails := range dataset2{
		for i := range emails{
			charSum += len(emails[i])
		}
		numEmails += len(emails)
	}
	assert.Equal(t, charSum/numEmails >= 21, true)
	assert.Equal(t, charSum/numEmails <= 23, true)
}

func TestUniqueEmails(t *testing.T){
	// test to make sure that millions of emails produced
	// by the DataGen func is all unique

	// using a GO map as a set (keys are unique)
	temp := make(map[string]bool)

	var totalEmails int

	data := GenData(10, 100000, 200000)

	for _, emails := range data{
		for i := range emails{
			temp[emails[i]] = true
		}
		totalEmails += len(emails)
	}

	// If the number of unique keys is equal to
	// the number of total emails generated
	// then all emails are unique
	assert.Equal(t, totalEmails, len(temp))
}
