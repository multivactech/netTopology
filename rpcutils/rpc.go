package rpcutils

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// PeersInfo 仅用于 rpc 通讯
type PeersInfo [][]string

// QueryPeersInfo 通过rpc，发送指定节点的所有peer信息，并返回，如果失败，返回空，和错误信息。
func QueryPeersInfo(ip string, rpcPort int) ([]string, []string, error) {
	clientPath, err := GetExecCmd()
	if err != nil {
		return nil, nil, err
	}
	url := fmt.Sprintf("%v:%v", ip, rpcPort)

	result := PeersInfo{}
	cmd := exec.Command(clientPath, "--notls", "-s", url, "-u", "mtvac", "-P", "mtvac", "getoutboundpeers")
	// log.Printf("%v --notls -s %v -u mtvac -P mtvac getoutboundpeers", clientPath, url)
	out, err := cmd.Output()
	if err != nil {
		log.Print("Got error when query shard info: ", cmd.Args, " ", err.Error())
		return nil, nil, err
	}
	json.Unmarshal(out, &result)
	// log.Print(result)
	return result[0], result[1], nil
}

//GetExecCmd get the executable binary testnet
func GetExecCmd() (string, error) {
	clientPath, err := exec.LookPath("testnet")
	if err != nil {
		log.Print("Failed to find testnet executable binary: ", err.Error())
		return "", err
	}
	return clientPath, nil
}
