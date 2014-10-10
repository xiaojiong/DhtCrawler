package DhtCrawler

import (
	"io"
	"log"
)

type DhtNode struct {
	node    *KNode
	table   *KTable
	network *Network
	log     *log.Logger
	master  chan string
	krpc    *KRPC
	dao     *Dao
}

func NewDhtNode(id *Id, logger io.Writer, dao *Dao, master chan string) *DhtNode {
	node := new(KNode)
	node.Id = *id

	dht := new(DhtNode)
	dht.dao = dao
	dht.log = log.New(logger, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	dht.node = node
	dht.table = new(KTable)
	dht.network = NewNetwork(dht)
	dht.krpc = NewKRPC(dht)
	dht.master = master

	return dht
}

func (dht *DhtNode) Run() {

	go func() { dht.network.Listening() }()
	go func() { dht.NodeFinder() }()

	for {
		select {
		case msg := <-dht.master:
			dht.log.Println(msg)
		}
	}
}
