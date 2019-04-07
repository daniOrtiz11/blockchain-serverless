package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	ma "gx/ipfs/QmTZBfrPJmjWsCvHEtX5FE6KimVJhsJg5sBbqEFYf4UZtL/go-multiaddr"
	peer "gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
	pstore "gx/ipfs/QmaCTz9RkrU13bm9kMB54f7atgqM4qkjDZpRwRoJiWXEqs/go-libp2p-peerstore"
	golog "gx/ipfs/QmbkT7eMTyXfpeyB3ZMxxcxg7XH8t6uXp49jqzz4HB7BGF/go-log"
	gologging "gx/ipfs/QmcaSwFc5RBg8yCq54QURwEU4nwjfCpjbpmaAm4VbdGLKv/go-logging"

	"github.com/davecgh/go-spew/spew"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index       int
	Timestamp   string
	CustomValue int
	Hash        string
	PrevHash    string
}

var cmdConsole string

// Blockchain is a series of validated Blocks
var Blockchain []Block

var mutex = &sync.Mutex{}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil || str == "" {
			return
		} else if str != "\n" {
			chain := make([]Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Printf(exceptionJSON)
				log.Fatal(err)
			}
			mutex.Lock()
			//Updating blc from broadcast
			Blockchain = updateBlc(chain, Blockchain)
			mutex.Unlock()
		}
	}
}

func writeData(rw *bufio.ReadWriter) {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			closeCon()
			log.Fatalln(sig)
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Printf(exceptionJSON)
				log.Println(err)
			}
			mutex.Unlock()
			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()

		}
	}()

	stdReader := bufio.NewReader(os.Stdin)

	for {
		inOk := false
		for inOk == false {
			showHelp()
			sendData, err := stdReader.ReadString('\n')
			if err != nil {
				log.Printf(exceptionReader)
				log.Fatal(err)
			}
			sendData = strings.Replace(sendData, "\n", "", -1)
			option, err := strconv.Atoi(sendData)
			if err == nil {
				switch option {
				case 1:
					viewState(rw)
					inOk = true
				case 2:
					insertBlock()
					inOk = true
				case 3:
					showRunCmd()
					inOk = true
				case 4:
					closeCon()
					inOk = true
				default:
					fmt.Println(badFormatOption)
				}
			} else {
				fmt.Println(badFormatNumber)
			}
		}
	}
}

func main() {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), ""}

	Blockchain = append(Blockchain, genesisBlock)

	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info
	//golog.SetAllLoggers(gologging.DEBUG) // Change to DEBUG for extra info

	// Parse options from the command line
	listenF := flag.Int(flagL, 0, "")
	target := flag.String(flagD, "", "")
	seed := flag.Int64(flagSeed, 0, "")
	mode := flag.String(flagMode, "", "")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal(badFormatArgument)
	}

	if *mode == debugstr {
		debug = true
	} else if *mode == prodstr || *mode == "" {
		debug = false
	} else {
		log.Fatal(badFormatMode)
	}

	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(*listenF, *seed)
	if err != nil {
		log.Fatal(err)
	}

	if *target == "" {
		log.Println(cmdInitialNode)
		log.Println(cmdInitialNode2)
		println(cmdConsole)
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler(p2p, handleStream)
		select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		ha.SetStreamHandler(p2p, handleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(*target)
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

		log.Println(cmdClientNode)
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, p2p)
		if err != nil {
			log.Fatalln(err)
		}
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)

		select {} // hang forever

	}
}

func showHelp() {
	fmt.Println(optionsTitle)
	fmt.Println(options1)
	fmt.Println(options2)
	fmt.Println(options3)
	fmt.Println(options4)
	fmt.Print("> ")
}

func viewState(rw *bufio.ReadWriter) {
	fmt.Println(options1Title)
	bytes, err := json.Marshal(Blockchain)
	if err != nil {
		log.Println(err)
	}
	spew.Dump(Blockchain)
	mutex.Lock()
	rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
	rw.Flush()
	mutex.Unlock()
}

func closeCon() {
	log.Fatal(endMessage)
}

func showRunCmd() {
	println(cmdConsole)
}

func insertBlock() {
	fmt.Println(options2Title)
	fmt.Print("> ")
	stdReader := bufio.NewReader(os.Stdin)
	sendData, err := stdReader.ReadString('\n')
	if err != nil {
		log.Printf(exceptionReader)
		log.Fatal(err)
	}
	sendData = strings.Replace(sendData, "\n", "", -1)
	customvalue, err := strconv.Atoi(sendData)
	if err == nil {
		Blockchain = insertBlc(customvalue, Blockchain)
	}
}
