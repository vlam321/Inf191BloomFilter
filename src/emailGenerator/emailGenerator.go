// Script that mainly for generating random email address strings
/*
Basic Requirements:
	- Function must be able to take an int input and generate
	that amount of "unique" email addresses
	- Average length of an email address should be between
	21-23 chars

*/
package emailGenerator

import (
	"math/rand"
	"time"
)

var seededRandPtr * rand.Rand = rand.New(rand.NewSource(
					time.Now().UnixNano()))

func genRandAdd() (string){
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

func GenEmailAddrs(size int) ([]string){
	emailAdds := make([]string, size)
	for i := range emailAdds{
		emailAdds[i] = genRandAdd()
	}
	return emailAdds
}
