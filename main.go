package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	ma "gx/ipfs/QmTZBfrPJmjWsCvHEtX5FE6KimVJhsJg5sBbqEFYf4UZtL/go-multiaddr"
	peer "gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
	pstore "gx/ipfs/QmaCTz9RkrU13bm9kMB54f7atgqM4qkjDZpRwRoJiWXEqs/go-libp2p-peerstore"
	golog "gx/ipfs/QmbkT7eMTyXfpeyB3ZMxxcxg7XH8t6uXp49jqzz4HB7BGF/go-log"
	gologging "gx/ipfs/QmcaSwFc5RBg8yCq54QURwEU4nwjfCpjbpmaAm4VbdGLKv/go-logging"
)

// Block represents each 'item' in the Blockchain
type Block struct {
	Index       int
	Timestamp   string
	Transaction Transaction
	Hash        string
	PrevHash    string
}

// Account struct
type Account struct {
	PublicID  string
	PrivateID string
	Name      string
	Amount    int
}

//Transaction struct
type Transaction struct {
	SourceID string
	TargetID string
	Amount   int
}

//LocalP2P host ip, port and key
type LocalP2P struct {
	Ipdir   string
	Port    string
	Key     string
	PrevKey string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

// Bank is a series of Accounts
var Bank []Account

//local p2p dir
var localP2P LocalP2P

//local account
var account Account

var mutex = &sync.Mutex{}

var targetP2P string

var listenF *int
var seed *int64

var logged bool

func main() {
	t := time.Now()
	genesisBlock := Block{}
	var transaction Transaction
	transaction.Amount = initAmount
	transaction.SourceID = initSource
	transaction.TargetID = initTarget

	genesisBlock = Block{0, t.String(), transaction, calculateHash(genesisBlock), ""}

	Blockchain = append(Blockchain, genesisBlock)

	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info
	//golog.SetAllLoggers(gologging.DEBUG) // Change to DEBUG for extra info

	// Parse options from the command line
	listenF = flag.Int(flagL, 0, "")
	seed = flag.Int64(flagSeed, 0, "")
	flag.Parse()

	if *listenF == 0 {
		*listenF = defaultPort
		fmt.Println(defaultPortStr)
	}
	generalMain()
}

func generalMain() {
	logged = false
	targetP2P = getTargetP2P()
	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(*listenF, *seed)
	if err != nil {
		log.Fatal(err)
	}

	if targetP2P == "" {
		fmt.Println(cmdInitialNode)
		fmt.Println(cmdInitialNode2)
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				closeCon()
				log.Fatalln(sig)
			}
		}()

		ha.SetStreamHandler(p2p, handleStream)
		fmt.Println(startingSetP2P)
		ping := getPingP2P()
		if ping != okC {
			setTargetP2P()
		} else {
			fmt.Println(errorDirRepeat)
			log.Fatal(helpRunning)

		}
		select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		ha.SetStreamHandler(p2p, handleStream)
		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(targetP2P)
		if err != nil {
			log.Fatalln(err)
		}
		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf(ipfs, peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		fmt.Println(cmdClientNode)
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, p2p)
		if err != nil {
			fmt.Println(errorDirRepeat)
			log.Fatal(helpRunning)
		}
		fmt.Println(startingSetP2P)
		ping := getPingP2P()
		if ping != okC {
			setTargetP2P()
		} else {
			fmt.Println(errorDirRepeat)
			log.Fatal(helpRunning)
		}
		fmt.Println(startedSetP2P)
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)

		select {} // hang forever

	}
}
