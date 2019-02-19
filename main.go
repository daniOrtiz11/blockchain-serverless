package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

/*
godotoenv-> get params from configuration files
spew -> print go structures
bufio -> buffer to read data
io -> handle out console
net -> network interface (TCP currently)
os -> read from file
strconv -> string conversion
sync -> dependencie to use mutex
*/

// Custom Block
type Block struct {
	Index         int
	Timestamp     string
	Private_value int
	Hash          string
	PrevHash      string
}

// Array Blocks -> real Blockchain
var Blockchain []Block

//channel for handle current access to Blockchain
var bcServer chan []Block

//mutex needed in sumalted broadcasting
var mutex = &sync.Mutex{}

func main() {
	//check access configuration file
	/*err := godotenv.Load("configuration.env")
	if err != nil {
		log.Fatal(err)
	}*/

	bcServer = make(chan []Block)

	// create genesis block
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	//get port from constant file
	httpPort := PORT

	// start TCP server
	server, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP Server Listening on port :", httpPort)
	defer server.Close()

	//Infinite loop to accetp and handle cons
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {

	defer conn.Close()

	io.WriteString(conn, "Enter a new Private_value:")

	scanner := bufio.NewScanner(conn)

	// take in Private_value from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			Private_value, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], Private_value)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

			bcServer <- Blockchain
			io.WriteString(conn, "\nEnter a new Private_value:")
		}
	}()

	// simulate receiving broadcast
	go func() {
		for {
			time.Sleep(30 * time.Second)
			mutex.Lock()
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range bcServer {
		spew.Dump(Blockchain)
	}

}

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// make sure the chain we're checking is longer than the current blockchain
func replaceChain(newBlocks []Block) {
	mutex.Lock()
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
	mutex.Unlock()
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, Private_value int) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Private_value = Private_value
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)

	return newBlock, nil
}
