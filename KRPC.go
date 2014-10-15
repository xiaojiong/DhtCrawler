package DhtCrawler

import (
	"bytes"
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
		var ok bool
		message := new(KRPCMessage)

		message.T, ok = val["t"].(string) //请求tid
		if !ok {
			return nil
		}

		message.Y, ok = val["y"].(string) //请求类型
		if !ok {
			return nil
		}

		message.Addr = raddr

		switch message.Y {
		case "q":
			query := new(Query)
			query.Y = val["q"].(string)
			query.A = val["a"].(map[string]interface{})
			message.Addion = query
			break
		case "r":
			res := new(Response)
			res.R = val["r"].(map[string]interface{})
			message.Addion = res
			break
		default:
			krpc.Dht.log.Println("invalid message")
			break
		}

		switch message.Y {
		case "q":
			krpc.Query(message)
			break
		case "r":
			krpc.Response(message)
			break
		}

	}
	return nil
}

func (krpc *KRPC) Response(msg *KRPCMessage) {
	if response, ok := msg.Addion.(*Response); ok {
		if nodestr, ok := response.R["nodes"].(string); ok {
			nodes := ParseBytesStream([]byte(nodestr))
			for _, node := range nodes {
				krpc.Dht.table.Put(node)
			}
		}
	}
}

func (krpc *KRPC) Query(msg *KRPCMessage) {
	if query, ok := msg.Addion.(*Query); ok {
		//krpc.Dht.log.Println(query.Y) 干掉输出避免磁盘写盘的情况
		if query.Y == "get_peers" {

			if infohash, ok := query.A["info_hash"].(string); ok {
				//krpc.Dht.log.Println(Id(infohash).String())

				krpc.Dht.dao.HashIns1.Exec(Id(infohash).String())
				//ih := Id(infohash)
				krpc.Dht.msq.addMessage(Id(infohash).String(), 1)

				//	krpc.Dht.log.Println(Id(infohash).String())
				nodes := ConvertByteStream(krpc.Dht.table.Snodes)
				data, _ := krpc.EncodingNodeResult(msg.T, "asdf13e", nodes)
				krpc.Dht.network.Send([]byte(data), msg.Addr)
			}
		}

		if query.Y == "announce_peer" {
			if infohash, ok := query.A["info_hash"].(string); ok {
				//krpc.Dht.dao.HashIns2.Exec(Id(infohash).String())
				krpc.Dht.msq.addMessage(Id(infohash).String(), 2)
			}
		}
	}
}

func ConvertByteStream(nodes []*KNode) []byte {
	buf := bytes.NewBuffer(nil)
	for _, v := range nodes {
		convertNodeInfo(buf, v)
	}
	return buf.Bytes()
}

func convertNodeInfo(buf *bytes.Buffer, v *KNode) {
	buf.Write(v.Id)
	convertIPPort(buf, v.Ip, v.Port)
}
func convertIPPort(buf *bytes.Buffer, ip net.IP, port int) {
	buf.Write(ip.To4())
	buf.WriteByte(byte((port & 0xFF00) >> 8))
	buf.WriteByte(byte(port & 0xFF))
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

type KRPCMessage struct {
	T      string
	Y      string
	Addion interface{}
	Addr   *net.UDPAddr
}

type Query struct {
	Y string
	A map[string]interface{}
}

type Response struct {
	R map[string]interface{}
}

func (krpc *KRPC) EncodingNodeResult(tid string, token string, nodes []byte) (string, error) {
	v := make(map[string]interface{})
	v["t"] = tid
	v["y"] = "r"
	args := make(map[string]string)
	args["id"] = string(krpc.Dht.node.Id)
	if token != "" {
		args["token"] = token
	}
	args["nodes"] = bytes.NewBuffer(nodes).String()
	v["r"] = args
	//krpc.Dht.log.Println(v)
	s, err := bencode.EncodeString(v)
	return s, err
}
