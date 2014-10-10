package godht

import (
	"github.com/zeebo/bencode"
	"math"
	"net"
	"sync/atomic"
)

type action func(arg map[string]interface{}, raddr *net.UDPAddr)

type KRPC struct {
	Dht   *DhtNode
	Types map[string]action
	tid   uint32
}

func NewKRPC(dhtNode *DhtNode) *KRPC {
	krpc := new(KRPC)
	krpc.Dht = dhtNode

	return krpc
}

func (krpc *KRPC) GenTID() uint32 {
	return krpc.autoID() % math.MaxUint16
}

func (krpc *KRPC) autoID() uint32 {
	return atomic.AddUint32(&krpc.tid, 1)
}

func (krpc *KRPC) Decode(data string, raddr *net.UDPAddr) error {
	val := make(map[string]interface{})

	if err := bencode.DecodeString(data, &val); err != nil {
		return err
	} else {
		msgType, ok := val["y"].(string) //请求类型
		if !ok {
			return nil
		}

		if msgType == "r" {
			krpc.Response(val, raddr)
		}
		if msgType == "q" {
			krpc.Query(val, raddr)
		}
	}
	return nil
}

func (krpc *KRPC) Response(arg map[string]interface{}, raddr *net.UDPAddr) {
	msgType, ok := arg["r"].(map[string]interface{})
	if ok && msgType["nodes"] != "" {
		if nodestr, ok := msgType["nodes"].(string); ok {
			nodes := ParseBytesStream([]byte(nodestr))
			for _, node := range nodes {
				krpc.Dht.table.Put(node)
			}
		}
	}
}
func (krpc *KRPC) Query(arg map[string]interface{}, raddr *net.UDPAddr) {
	krpc.Dht.log.Println(arg["q"])
}

func ParseBytesStream(data []byte) []*KNode {
	var nodes []*KNode = nil
	for j := 0; j < len(data); j = j + 26 {
		if j+26 > len(data) {
			break
		}

		kn := data[j : j+26]
		node := new(KNode)
		node.Id = Id(kn[0:20])
		node.Ip = kn[20:24]
		port := kn[24:26]
		node.Port = int(port[0])<<8 + int(port[1])
		nodes = append(nodes, node)
	}
	return nodes
}
