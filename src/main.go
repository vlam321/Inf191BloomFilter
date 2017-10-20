package main

import "./emailGenerator"
import "fmt"

func getAverageLen(addrs []string)(int){
	numEle := len(addrs)
	numChar := 0
	for i := range addrs{
		numChar += len([]rune(addrs[i]))
	}
	return numChar/numEle
}

func main(){
	fmt.Println(getAverageLen(emailGenerator.GenEmailAddrs(100)))
	fmt.Println(getAverageLen(emailGenerator.GenEmailAddrs(1000)))
	fmt.Println(getAverageLen(emailGenerator.GenEmailAddrs(10000)))
	fmt.Println(getAverageLen(emailGenerator.GenEmailAddrs(100000)))
	fmt.Println(getAverageLen(emailGenerator.GenEmailAddrs(1000000)))
}
