package main

import (
	"fmt"
	"github.com/xiaojiong/DhtCrawler"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	master := make(chan string)

	msq := DhtCrawler.NewMSQ("127.0.0.1:123456")
	dao := DhtCrawler.NewDao("user", "password", "127.0.0.1", 3306, "test")
	//进程数量
	for i := 0; i < 20; i++ {
		go func() {
			id := DhtCrawler.GenerateID()
			dhtNode := DhtCrawler.NewDhtNode(&id, os.Stdout, dao, msq, master)

			dhtNode.Run()
		}()
	}

	for {
		select {
		case msg := <-master:
			fmt.Println(msg)
		}
	}
}
