package godht

import (
	"net"
)

type KTable struct {
	Nodes []*KNode
}

func (table *KTable) Put(node *KNode) {
	table.Nodes = append(table.Nodes, node)
}

type KNode struct {
	Id   Id
	Ip   net.IP
	Port int
}
