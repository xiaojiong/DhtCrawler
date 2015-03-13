package DhtCrawler

import (
	"net"
)

type Network struct {
	Dht  *DhtNode
	Conn *net.UDPConn
}

func NewNetwork(dhtNode *DhtNode) *Network {
	nw := new(Network)
	nw.Dht = dhtNode
	nw.Init()
	return nw
}
func (nw *Network) Init() {
	addr := new(net.UDPAddr)

	var err error
	nw.Conn, err = net.ListenUDP("udp", addr)

	if err != nil {
		panic(err)
	}

	laddr := nw.Conn.LocalAddr().(*net.UDPAddr)
	nw.Dht.node.Ip = laddr.IP
	nw.Dht.node.Port = laddr.Port
}

func (nw *Network) Listening() {
	buf := make([]byte, 1000)
	for {
		_, raddr, err := nw.Conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		nw.Dht.krpc.Decode(string(buf), raddr)
	}
}

func (nw *Network) Send(m []byte, addr *net.UDPAddr) error {
	_, err := nw.Conn.WriteToUDP(m, addr)

	if err != nil {
		nw.Dht.log.Println(err)
	}
	return err
}
