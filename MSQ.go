package DhtCrawler

import (
	"fmt"
	"github.com/xiaojiong/memcache"
	"time"
)

type MSQ struct {
	mem *memcache.Connection
}

func NewMSQ(dns string) *MSQ {
	mq := new(MSQ)

	mc, err := memcache.Connect(dns, 2*time.Second)

	if err != nil {
		fmt.Printf("mq error Connect: %v\n", err)
	}

	mq.mem = mc
	return mq
}

func (mq *MSQ) addMessage(hash string, queueIdx int) {

	stored, err := mq.mem.Set(fmt.Sprintf("dht-hash-%d", queueIdx), 0, 0, []byte(hash))
	if err != nil {
		fmt.Printf("mq error Set: %v\n", err)
	}
	if !stored {
		fmt.Printf("mq error want true, got %v\n", stored)
	}
}
