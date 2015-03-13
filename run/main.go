package main

import (
	"DhtCrawler"
	"fmt"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(2)

	//主进程
	master := make(chan string)

	//爬虫输出抓取到的hashIds通道
	outHashIdChan := make(chan string)

	//开启的dht节点
	for i := 0; i < 2; i++ {
		go func() {
			id := DhtCrawler.GenerateID()
			dhtNode := DhtCrawler.NewDhtNode(&id, os.Stdout, outHashIdChan, master)

			dhtNode.Run()
		}()
	}

	for {
		select {

		//输出爬虫抓取的HashId结果
		case hashId := <-outHashIdChan:
			fmt.Println(hashId)

		case msg := <-master:
			fmt.Println(msg)
		}
	}
}
