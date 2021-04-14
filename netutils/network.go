package netutils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/multivactech/netTopo/configutils"
	"github.com/multivactech/netTopo/rpcutils"
)

// Network 表示当前网络信息
type Network struct {
	Nodes         []*Node
	Mutex         *sync.Mutex
	OnlineNodeNum int
	Dic           map[string]int
}

// Init 网络初始化，加载网络初始节点
func (nw *Network) Init() {
	nw.Mutex = new(sync.Mutex)
	nw.Dic = make(map[string]int)
	nw.AddNode(configutils.FirstNodeIP, configutils.FirstNodeRPCPort)
}

// Update 更新网络，对每个node进行rpc请求，取得其peer信息
func (nw *Network) Update() {
	log.Print("开始更新网络信息")
	nw.Mutex.Lock()
	defer nw.Mutex.Unlock()

	allNodes := &Set{}

	for _, node := range nw.Nodes {
		// rpc 调用获得所有的 inbound 和 outbound
		outbounds, inbounds, err := rpcutils.QueryPeersInfo(node.IP, node.RPCPort)
		if err != nil {
			log.Print(err)
			// 请求失败，说明节点不可达，可能为掉线，更新状态
			node.Status = OFFLINE
			continue
		}
		// 将 inbound 和 outbound 都添加到新节点集合中
		allNodes.AddAll(inbounds, outbounds)
		log.Printf("allNodes num is : %v", allNodes.Size())

		// 更新节点的outbound
		node.UpdateOutbound(outbounds)
		// 根据inbound，找到from节点，然后更新outbounds
		nw.UpdateNodeInbounds(inbounds, node.GetRPCAddr())
	}
	log.Print("合并新节点到维护的节点列表中")

	allNodesInfo := allNodes.GetAll()
	for _, node := range allNodesInfo {
		strs := strings.Split(node, ":")
		ip := strs[0]
		rpcPort, _ := strconv.Atoi(strs[1])

		if !nw.IsExist(ip, rpcPort) {
			nw.AddNode(ip, rpcPort)
		} else {
			nw.UpdateNodeStatus(ip, rpcPort)
		}
	}
	// for _, node := range nw.Nodes {
	// 	log.Print(node.GetRPCAddr(), "节点outbound：")
	// 	node.OutBound.Print()
	// }
}

// GetAllNodeFromBound 从获取到的inbound和outbound中统计所有的节点
func (nw *Network) GetAllNodeFromBound(inbounds, outbounds []string) *Set {
	allNodes := &Set{}
	for _, peer := range inbounds {
		allNodes.Add(peer)
	}
	for _, peer := range outbounds {
		allNodes.Add(peer)
	}
	return allNodes
}

// AddNode 给当前网络增加一个节点，指定ip和port
func (nw *Network) AddNode(ip string, rpcPort int) error {
	index := len(nw.Nodes)
	nw.Nodes = append(nw.Nodes, &Node{
		IP:      ip,
		Port:    rpcPort - 1,
		RPCPort: rpcPort,
		Status:  ONLINE,
	})
	nw.Dic[fmt.Sprintf("%v:%v", ip, rpcPort)] = index
	return nil
}

// IsExist 判断当前维护的网络中是否存在该ip
func (nw *Network) IsExist(ip string, rpcPort int) bool {
	for _, node := range nw.Nodes {
		if node.IP == ip && node.RPCPort == rpcPort {
			return true
		}
	}
	return false
}

// UpdateNodeStatus 更新节点状态
func (nw *Network) UpdateNodeStatus(ip string, rpcPort int) {
	for _, node := range nw.Nodes {
		if node.IP == ip && node.RPCPort == rpcPort {
			node.Status = ONLINE
			return
		}
	}
}

// UpdateNodeInbounds 根据inbounds 更新节点的 outbounds
func (nw *Network) UpdateNodeInbounds(inbounds []string, out string) {
	for _, inbound := range inbounds {
		inboundIndex, ok := nw.Dic[inbound]
		if ok {
			nw.Nodes[inboundIndex].UpdateOutbound([]string{out})
		}
	}
}
