// Script that mainly for generating random email address strings
/*
Basic Requirements:
	- Function must be able to take an int input and generate
	that amount of "unique" email addresses
	- Average length of an email address should be between
	21-23 chars

*/
package bloomDataGenerator

import (
	"math/rand"
	"time"
)

const USER_ID_LIMIT = 50;

// random seed
var seededRandPtr * rand.Rand = rand.New(rand.NewSource(
					time.Now().UnixNano()))

func genRandAddr() (string){
	// creates one random email address 
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	userNameLen := seededRandPtr.Intn(10) + 5
	domainLen := seededRandPtr.Intn(5) + 6
	userNameArray := make([]byte, userNameLen)
	domainName := make([]byte, domainLen)

	for i := range userNameArray {
		userNameArray[i] = chars[seededRandPtr.Intn(len(chars))]
	}

	for i := range domainName {
		domainName[i] = chars[seededRandPtr.Intn(len(chars))]
	}

	return string(userNameArray) + "@" + string(domainName) + ".com"
}

func genEmailAddrs(min, max int) ([]string){
	// creates array[size] of email addresses
	size := seededRandPtr.Intn(max - min) + min
	emailAdds := make([]string, size)
	for i := range emailAdds{
		emailAdds[i] = genRandAddr()
	}
	return emailAdds
}

func genUsers(size int) ([]int){
	// create array[size] of user ids
	users := make([]int, size)
	for i := range users{
		users[i] = i
	}
	return users
}

func GenData(user_size, min_email_addrs, max_email_addrs int) (map[int][]string){
	// create map(int, string[]) (user_id, email addresses) 
	randData := make(map[int][]string)
	users := genUsers(user_size)
	for i := range users{
		randData[users[i]] = genEmailAddrs(min_email_addrs, max_email_addrs)
	}
	return randData
}

