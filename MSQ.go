package DhtCrawler

import (
	"fmt"
	"github.com/xiaojiong/memcache"
	"time"
)

type MSQ struct {
	mem      *memcache.Connection
	hashBox1 []*string
	hashBox2 []*string
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

	var sendVal string
	if queueIdx == 1 {
		mq.hashBox1 = append(mq.hashBox1, &hash)
		if len(mq.hashBox1) == 20 {
			for _, hash := range mq.hashBox1 {
				sendVal = sendVal + *hash + "{||}"
			}
			mq.hashBox1 = nil
		}
	}

	if queueIdx == 2 {
		mq.hashBox2 = append(mq.hashBox2, &hash)
		if len(mq.hashBox2) == 20 {
			for _, hash := range mq.hashBox2 {
				sendVal = sendVal + *hash + "{||}"
			}
			mq.hashBox2 = nil
		}
	}

	if sendVal != "" {
		stored, err := mq.mem.Set(fmt.Sprintf("dht-hash-%d", queueIdx), 0, 0, []byte(sendVal))
		if err != nil {
			fmt.Printf("mq error Set: %v\n", err)
		}
		if !stored {
			s, _ := mq.mem.Set(fmt.Sprintf("dht-hash-%d", queueIdx), 0, 0, []byte(sendVal))
			if !s {
				fmt.Printf("mq error want true, got %v\n", stored)
			}
		}
	}
}
