package main

import "fmt"
import "github.com/marpaia/graphite-golang"
import "log"
import "time"

/*type server struct {
	//foo
	//bar
	g *graphite.Graphite

}

//func Newserver() *server {
//	var s server
	
		g, err := graphite.NewGraphite("localhost", 2003)
		if err != nil {
			panic(err)
}*/
func main() {

	//var s server

	g, err := graphite.NewGraphite("localhost", 2003)
	if err != nil {
		panic(err)
	}
	for i:=0; i<20; i++ {
		time.Sleep(time.Second)
		err = g.SimpleSend("foo.bar.metric_one", fmt.Sprintf("%d", i))
		if err != nil {
			log.Printf("err: %v", err)
			return
		}
		err = g.SimpleSend("foo.bar.metric_two", fmt.Sprintf("%d", 20 - i))
		if err != nil {
			log.Printf("err: %v", err)
			return
		}
	}
	fmt.Println("hello world")
}
