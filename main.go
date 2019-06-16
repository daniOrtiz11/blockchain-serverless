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

// Block struct
type Block struct {
	Index       int
	Timestamp   string
	Transaction Transaction
	Hash        string
	PrevHash    string
}

//LocalP2P struct
type LocalP2P struct {
	Ipdir   string
	Port    string
	Key     string
	PrevKey string
}

//Blockchain is a series of Blocks
var Blockchain []Block

//Bank is a series of Accounts
var Bank []Account

//localP2P is the local dir node
var localP2P LocalP2P

//account is the account of user logged
var account Account

var mutex = &sync.Mutex{}

var targetP2P string

//port
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

	//initial block
	genesisBlock = Block{0, t.String(), transaction, calculateHash(genesisBlock), ""}

	Blockchain = append(Blockchain, genesisBlock)

	golog.SetAllLoggers(gologging.INFO)
	//golog.SetAllLoggers(gologging.DEBUG) // Change to DEBUG for extra info

	// Parse options from the command line
	listenF = flag.Int(flagL, 0, "")
	seed = flag.Int64(flagSeed, 0, "")
	flag.Parse()

	if *listenF == 0 {
		*listenF = defaultPort
		fmt.Println(defaultPortStr)
	}
	generalMain(false)
}

/*
Main Function
*/
func generalMain(reconnected bool) {
	logged = false

	//get dir target to connect
	targetP2P = getTargetP2P()
	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(*listenF, *seed)
	if err != nil {
		log.Fatal(err)
	}
	//Initial node
	if targetP2P == "" {
		restartLog()
		fmt.Println(cmdInitialNode)
		fmt.Println(cmdInitialNode2)
		logEntry(blockchainstr, "", 3)
		//Detect forced exit Ctrl + C
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				closeCon()
				log.Fatalln(sig)
			}
		}()

		//Opening streaming
		ha.SetStreamHandler(p2p, handleStream)
		fmt.Println(startingSetP2P)

		//check dir exists
		ping := getPingP2P()
		if ping != okC {
			setTargetP2P()
		} else {
			fmt.Println(errorDirRepeat)
			log.Fatal(helpRunning)

		}
		select {} // hang forever
	} else {
		//Opening streaming
		ha.SetStreamHandler(p2p, handleStream)

		//Get target's peer id from the given multiaddress
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

		//Decapsulate the /ipfs/<peerID> part from the target
		targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf(ipfs, peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		//LibP2P to contact peerid and targetAddr
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)
		fmt.Println(cmdClientNode)
		// make a new stream
		s, err := ha.NewStream(context.Background(), peerid, p2p)
		if err != nil {
			fmt.Println(errorDirRepeat)
			log.Fatal(helpRunning)
		}
		fmt.Println(startingSetP2P)
		ping := getPingP2P()
		if ping != okC {
			//Dir ok and go to set in server
			setTargetP2P()
		} else {
			fmt.Println(errorDirRepeat)
			log.Fatal(helpRunning)
		}
		fmt.Println(startedSetP2P)

		// Create a buffered stream to read and write between nodes
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write.
		go writeData(rw, reconnected)
		go readData(rw)

		select {} // hang forever

	}
}
