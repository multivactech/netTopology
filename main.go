package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/multivactech/netTopo/htmlutils"
	"github.com/multivactech/netTopo/netutils"
)

func main() {
	network := netutils.Network{}
	network.Init()

	http.HandleFunc("/", HandleGetTopo)
	go http.ListenAndServe(":8000", nil)

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	done := make(chan bool)

	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case <-ticker.C:
			network.Update()
			htmlutils.CreateHTML(network.Nodes)
		}
	}
}

// HandleGetTopo 返回拓扑图
func HandleGetTopo(w http.ResponseWriter, r *http.Request) {
	w.Write(getHTML())
}

func getHTML() []byte {
	op, _ := os.Open("htmlutils/topo.html")
	defer op.Close()
	data := make([]byte, 1000000)
	n, _ := op.Read(data)
	return data[0:n]
}
