package main

import (
	"DhtCrawler"
	"fmt"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(16)
	master := make(chan string)

	dao := DhtCrawler.NewDao("user", "password", "127.0.0.1", 3306, "test")

	for i := 0; i < 20; i++ {
		go func() {
			id := DhtCrawler.GenerateID()
			a := DhtCrawler.NewDhtNode(&id, os.Stdout, dao, master)

			a.Run()
		}()
	}

	for {
		select {
		case msg := <-master:
			fmt.Println(msg)
		}
	}
}
