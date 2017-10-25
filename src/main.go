package main

import "./bloomDataGenerator"
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
	m := bloomDataGenerator.GenData(5, 10, 100)
	for i := range m{
		fmt.Println(i, "---", m[i], "\n")	
	}
}

