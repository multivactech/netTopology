package rpcutils

import (
	"testing"
)

func TestQueryPeersInfo(t *testing.T) {
	ip := "13.251.185.134"
	port := 18334
	outbounds, inbounds, err := QueryPeersInfo(ip, port)
	t.Log(outbounds)
	t.Log(inbounds)
	t.Log(err)
}
