package netutils

import (
	"fmt"
)

const (
	// ONLINE is the node is online and can connected.
	ONLINE = iota
	// OFFLINE 节点掉线
	OFFLINE
)

// Node 表示网络中一个节点
type Node struct {
	IP      string
	Port    int
	RPCPort int
	Status  int
	// OutBound 当前节点连出的所有地址，"ip:RPCPort" 的形式
	OutBound Set
}

// UpdateOutbound 更新节点的outbound信息
func (node *Node) UpdateOutbound(peers []string) {
	for _, peer := range peers {

		// log.Printf("%v add outbound %v", node.GetRPCAddr(), peer)
		node.OutBound.Add(peer)
	}
}

// GetRPCAddr 获取当前节点的地址，格式为"ip:rpcPort"
func (node *Node) GetRPCAddr() string {
	return fmt.Sprintf("%v:%v", node.IP, node.RPCPort)
}

// GetAddr 获取当前节点的listen地址，格式为"ip:port"
func (node *Node) GetAddr() string {
	return fmt.Sprintf("%v:%v", node.IP, node.Port)
}
