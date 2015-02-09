// Runs a node on a random UDP port that attempts to collect 10 peers for an
// infohash, then keeps running as a passive DHT node.
//
// IMPORTANT: if the UDP port is not reachable from the public internet, you
// may see very few results.
//
// To collect 10 peers, it usually has to contact some 1k nodes. It's much easier
// to find peers for popular infohashes. This process is not instant and should
// take a minute or two, depending on your network connection.
//
//
// There is a builtin web server that can be used to collect debugging stats
// from http://localhost:8711/debug/vars.
package main

import (
	"flag"
	"fmt"
	"encoding/hex"
	"net/http"
	"os"
	"log"
	"net"
	"time"

	"github.com/qsun/dht"
)

var infoHashQueryChan chan string
var DHTMapping map[int]*dht.DHT

const (
	httpPortTCP = 8711
	numTarget   = 10
	// exampleIH   = "deca7a89a1dbdc4b213de1c0d5351e92582f31fb" // ubuntu-12.04.4-desktop-amd64.iso
)

/* dht.Logger */
type InfoHashLogger struct {
	DHTId int
}

func (l *InfoHashLogger) GetPeers(peer net.UDPAddr, id string, infoHash dht.InfoHash) {
	log.Println("Query: ", peer, "id: ", hex.EncodeToString([]byte(id)),
		"InfoHash: ", hex.EncodeToString([]byte(infoHash)))
}

func main() {
	flag.Parse()

	DHTMapping = make(map[int]*dht.DHT)

	initControl()

	node_count := 0
	
	for node_count < 30 {
		// Starts a DHT node with the default options. It picks a random UDP port. To change this, see dht.NewConfig.
		d, err := dht.New(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "New DHT error: %v", err)
			os.Exit(1)

		}

		d.Logger = &InfoHashLogger{node_count}

		go d.Run()
		go drainresults(d)

		DHTMapping[node_count] = d

		node_count = node_count + 1
	}

	// For debugging.
	go http.ListenAndServe(fmt.Sprintf(":%d", httpPortTCP), nil)


	for {
		infoHashQuery := <- infoHashQueryChan
		log.Println("Retrieve: ", infoHashQuery)

		for _, dhtNet := range DHTMapping {
			
			dhtNet.PeerRequest(hex.
		}
		
		time.Sleep(5 * time.Second)
	}
}

// drainresults loops, printing the address of nodes it has found.
func drainresults(n *dht.DHT) {
	count := 0
	fmt.Println("=========================== DHT")
	fmt.Println("Note that there are many bad nodes that reply to anything you ask.")
	fmt.Println("Peers found:")
	for r := range n.PeersRequestResults {
		for _, peers := range r {
			for _, x := range peers {
				fmt.Printf("%d: %v\n", count, dht.DecodePeerAddress(x))
				count++
				if count >= numTarget {
					os.Exit(0)
				}
			}
		}
	}
}
