package htmlutils

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/multivactech/netTopo/netutils"
)

var section1 = `
<!doctype html>
<html>
<head>
  <title>Network | Basic usage</title>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/vis/4.21.0/vis.min.js"></script>
  <link href="https://cdnjs.cloudflare.com/ajax/libs/vis/4.21.0/vis.min.css" rel="stylesheet" type="text/css" />
  <style type="text/css">
    #mynetwork {
      width: 1920px;
      height: 1080px;
      border: 1px solid lightgray;
    }
  </style>
</head>
<body>
<p>
  Create a simple network with some nodes and edges.
</p>
<div id="mynetwork"></div>
<script type="text/javascript">
  // create an array with nodes
  var nodes = new vis.DataSet([
`
var section3 = `
]);
// create an array with edges
var edges = new vis.DataSet([
`

var section5 = `
]);
// create a network
var container = document.getElementById('mynetwork');
var data = {
  nodes: nodes,
  edges: edges
};
var options = {
	"edges": {
	  "smooth": {
		"forceDirection": "none",
		"roundness": 0.1
	  }
	},
	"physics": {
	  "barnesHut": {
		"gravitationalConstant": -19800,
		"centralGravity": 2.65,
		"springLength": 75,
		"springConstant": 0,
		"damping": 0.11
	  },
	  "minVelocity": 0.75
	}
  }
var network = new vis.Network(container, data, options);
</script>
</body>
</html>
`

var tableModel = "<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>"

// CreateHTML 根据所有node创建html
func CreateHTML(nodes []*netutils.Node) {
	log.Printf("开始更新html, nodes num: %v", len(nodes))
	section2 := ""
	dic := make(map[string]int)
	inboundDic := make(map[string]int)
	index := 1
	// 增加节点以及编号
	for _, node := range nodes {
		if node.Status == netutils.ONLINE && node.Port != 8333 && node.Port != 18333 {
			addr := node.GetRPCAddr()
			dic[addr] = index

			// section2 += fmt.Sprintf("    {id: %v, label: '%v'},\r\n", index, addr)
			inboundDic[addr] = 0
			index++
		}
	}

	section4 := ""

	for _, node := range nodes {
		from, ok := dic[node.GetRPCAddr()]
		if ok && node.Status == netutils.ONLINE {
			outbounds := node.OutBound.GetAll()
			log.Print("outbound size: ", node.OutBound.Size())
			log.Print(outbounds)
			for _, toNode := range outbounds {
				toIndex, ok := dic[toNode]
				if ok {
					section4 += fmt.Sprintf("    {from: %v, to: %v, arrows: 'to'},\r\n", from, toIndex)
					inboundDic[toNode]++
				}
			}
		}
	}
	for rpcAddr, nodeIndex := range dic {

		node := nodes[nodeIndex]
		if node.Status == netutils.ONLINE {

			section2 += fmt.Sprintf("    {id: %v, label: '%v, in: %v, out: %v'},\r\n", nodeIndex, rpcAddr, inboundDic[rpcAddr], node.OutBound.Size())
		}
	}

	message := []byte(section1 + section2 + section3 + section4 + section5)
	err := ioutil.WriteFile("htmlutils/topo.html", message, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
