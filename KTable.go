package DhtCrawler

import (
	"net"
)

type KTable struct {
	Nodes  []*KNode
	Snodes []*KNode
}

func (table *KTable) Put(node *KNode) {
	table.Nodes = append(table.Nodes, node)
	if len(table.Snodes) < 8 {
		table.Snodes = append(table.Snodes, node)
	}

}

type KNode struct {
	Id   Id
	Ip   net.IP
	Port int
}
