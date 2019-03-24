package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	crypto "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"
	ma "gx/ipfs/QmTZBfrPJmjWsCvHEtX5FE6KimVJhsJg5sBbqEFYf4UZtL/go-multiaddr"
	net "gx/ipfs/QmY3ArotKMKaL7YGfbQfyDrib6RVraLqZYWXZvVgZktBxp/go-libp2p-net"
	peer "gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
	host "gx/ipfs/QmYrWiWM4qtrnCeT3R14jY3ZZyirDNJgwK57q4qFYePgbd/go-libp2p-host"
	pstore "gx/ipfs/QmaCTz9RkrU13bm9kMB54f7atgqM4qkjDZpRwRoJiWXEqs/go-libp2p-peerstore"
	golog "gx/ipfs/QmbkT7eMTyXfpeyB3ZMxxcxg7XH8t6uXp49jqzz4HB7BGF/go-log"
	gologging "gx/ipfs/QmcaSwFc5RBg8yCq54QURwEU4nwjfCpjbpmaAm4VbdGLKv/go-logging"

	"github.com/davecgh/go-spew/spew"
	libp2p "github.com/libp2p/go-libp2p"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index       int
	Timestamp   string
	CustomValue int
	Hash        string
	PrevHash    string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

var mutex = &sync.Mutex{}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress.
func makeBasicHost(listenPort int, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	/*
		Seed to create a random address for our host.
	*/
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	/*
		opts: address, identity to new connect between peer.
	*/

	addrstr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)
	//addrstr := fmt.Sprintf("/ip4/192.168.1.135/tcp/%d", listenPort)
	//addrstr := fmt.Sprintf("/ip4/81.0.3.81/tcp/%d", listenPort)
	//addrstr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)
	//addrstr := fmt.Sprintf("/ip4/3.17.190.106/tcp/%d", listenPort)

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(addrstr),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addrs := basicHost.Addrs()
	var addr ma.Multiaddr
	findaddr := false
	// select the address starting with "ip4"
	for _, i := range addrs {
		if strings.HasPrefix(i.String(), "/ip4") {
			addr = i
			log.Printf("I am in for %s\n", i)
			findaddr = true
			//break
		}
	}
	/*
		Handle clase if there are any problem get addr:
	*/
	if findaddr == false {
		log.Printf("In if")
		addrcero := basicHost.Addrs()[0]
		fullAddrcero := addrcero.Encapsulate(hostAddr)
		sfulladdr := fmt.Sprintf("%s", fullAddrcero)
		s := strings.Split(sfulladdr, "/")
		nextkey := s[3]
		if nextkey == "tcp" {
			nextkey = s[6]
		}
		addrstr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)
		realaddress := addrstr + "/ipfs/" + nextkey
		log.Printf("I am %s\n", fullAddrcero)
		if debug {
			log.Printf("Now run \"go run main.go bcfunctions.go constants.go aws.go utils.go -l %d -d %s\" on a different terminal\n", listenPort+1, realaddress)
		} else {
			log.Printf("Now run \"blockchain-serverless -l %d -d %s\" on a different terminal\n", listenPort+1, realaddress)
		}
	} else {
		fullAddr := addr.Encapsulate(hostAddr)
		log.Printf("I am %s\n", fullAddr)
		if debug {
			log.Printf("Now run \"go run main.go bcfunctions.go constants.go aws.go utils.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
		} else {
			log.Printf("Now run \"blockchain-serverless -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
		}
	}

	return basicHost, nil
}

func handleStream(s net.Stream) {

	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil || str == "" {
			log.Printf("client lost")
			return
		} else if str != "\n" {
			chain := make([]Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Printf("Exception json.Unmarshal: ")
				log.Fatal(err)
			}
			mutex.Lock()
			if len(chain) > len(Blockchain) {
				Blockchain = chain
				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
				if err != nil {
					log.Printf("Exception json.MarshalIndent: ")
					log.Fatal(err)
				}
				if !debug {
					updateGlobal(bytes)
				}

				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			mutex.Unlock()
		}
	}
}

func writeData(rw *bufio.ReadWriter) {

	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Printf("Exception json.Marshal: ")
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
				log.Printf("Exception stdReader:")
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
					closeCon()
					inOk = true
				default:
					fmt.Println("Error: Please provide a option")
				}
			} else {
				fmt.Println("Error: Please provide a number")
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
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	mode := flag.String("mode", "", "mode debug")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	if *mode == "debug" {
		debug = true
	} else if *mode == "" {
		debug = false
	} else {
		log.Fatal("Please provide a correct mode with -mode")
	}

	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(*listenF, *seed)
	if err != nil {
		log.Fatal(err)
	}

	if *target == "" {
		log.Println("Initial node")
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)
		select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

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
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
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

func updateGlobal(bytes []byte) {
	bytestofile(bytes)
	uploadfile()
}

func showHelp() {
	fmt.Println("Options:")
	fmt.Println("1. View State")
	fmt.Println("2. Insert new value")
	fmt.Println("3. Close connection")
	fmt.Print("> ")
}

func viewState(rw *bufio.ReadWriter) {
	fmt.Println("Current Blockchain state:")
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
	fmt.Println("Here in option 3")
}

func insertBlock() {
	fmt.Println("Insert a new value:")
	fmt.Print("> ")
	stdReader := bufio.NewReader(os.Stdin)
	sendData, err := stdReader.ReadString('\n')
	if err != nil {
		log.Printf("Exception stdReader:")
		log.Fatal(err)
	}
	sendData = strings.Replace(sendData, "\n", "", -1)
	customvalue, err := strconv.Atoi(sendData)
	if err == nil {
		newBlock := generateBlock(Blockchain[len(Blockchain)-1], customvalue)

		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			mutex.Lock()
			Blockchain = append(Blockchain, newBlock)
			mutex.Unlock()
		}
	}
}
