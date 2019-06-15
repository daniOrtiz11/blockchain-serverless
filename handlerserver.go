package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	crypto "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"
	ma "gx/ipfs/QmTZBfrPJmjWsCvHEtX5FE6KimVJhsJg5sBbqEFYf4UZtL/go-multiaddr"
	net "gx/ipfs/QmY3ArotKMKaL7YGfbQfyDrib6RVraLqZYWXZvVgZktBxp/go-libp2p-net"
	host "gx/ipfs/QmYrWiWM4qtrnCeT3R14jY3ZZyirDNJgwK57q4qFYePgbd/go-libp2p-host"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
)

/*
Handle streaming
*/
func handleStream(s net.Stream) {

	// Create a buffered stream to read and write between nodes
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write.
	go writeData(rw)
	go readData(rw)

}

/*
Creates a LibP2P host with a random peer ID listening
*/
func makeBasicHost(listenPort int, randseed int64) (host.Host, error) {
	var r io.Reader
	//Seed to create a random address for our host.
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	//Generate a key pair for this host
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	//Using local ip and provided port
	ipLocal, err := getExternalIP()
	if err != nil {
		fmt.Println(err)
		closeCon()
	}
	addrstr := fmt.Sprintf(ipLocal, listenPort)

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(addrstr),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf(ipfs, basicHost.ID().Pretty()))
	addrs := basicHost.Addrs()
	var addr ma.Multiaddr
	findaddr := false
	// select the address starting with "ip4"
	for _, i := range addrs {
		if strings.HasPrefix(i.String(), ipv4) {
			addr = i
			findaddr = true
			//break
		}
	}
	//Handle if there are any problem getting addr
	if findaddr == false {
		addrcero := basicHost.Addrs()[0]
		fullAddrcero := addrcero.Encapsulate(hostAddr)
		sfulladdr := fmt.Sprintf("%s", fullAddrcero)
		s := strings.Split(sfulladdr, "/")
		nextkey := s[3]
		if nextkey == tcp {
			nextkey = s[6]
		}
		//Using local ip and provided port
		ipLocal, err := getExternalIP()
		if err != nil {
			fmt.Println(err)
			closeCon()
		}
		addrstr := fmt.Sprintf(ipLocal, listenPort)
		realaddress := addrstr + ipfs2 + nextkey
		//set node addr
		parserLocalP2P(realaddress)
	} else {
		fullAddr := addr.Encapsulate(hostAddr)
		//set node addr
		parserLocalP2P(fullAddr.String())
	}

	return basicHost, nil
}

/*
Get target from AWS
*/
func getTargetP2P() string {
	target := generalLambda(arnFuncGetP2P, "")
	targetParser := ""
	if target != "" && target != emptyC {
		targetParser = parserTarget(target)
	}
	return targetParser
}

/*
Check if addr exists in AWS
*/
func getPingP2P() string {
	clientP2P := localP2P.Ipdir + "/" + localP2P.Port + "/" + localP2P.Key + "/" + localP2P.PrevKey
	paramClientP2P := initParam + clientP2P + "\"}"
	ping := generalLambda(arnFuncPingP2P, paramClientP2P)
	pingResult := koC
	if ping == okC {
		pingResult = okC
	}
	return pingResult
}

/*
Set addr in AWS
*/
func setTargetP2P() {
	clientP2P := localP2P.Ipdir + "/" + localP2P.Port + "/" + localP2P.Key + "/" + localP2P.PrevKey
	paramClientP2P := initParam + clientP2P + "\"}"
	resp := generalLambda(arnFuncSetP2P, paramClientP2P)
	if resp == okC {
		setTargetP2P()
	} else if resp == koC {
		log.Fatal(errorSetP2P)
	}
}

/*
Set log in AWS
*/
func setLog(newentry string) {
	paramLog := newentry
	paramLogAws := initParam + paramLog + "\"}"
	resp := generalLambda(arnFuncLog, paramLogAws)
	if resp == okC {
		setLog(newentry)
	} else if resp == koC {
		log.Fatal(errorSetLog)
	}
}

/*
Delete addr in AWS
*/
func deleteTargetP2P(needCheck bool) {
	clientP2P := localP2P.Ipdir + "/" + localP2P.Port + "/" + localP2P.Key + "/" + localP2P.PrevKey
	paramClientP2P := initParam + clientP2P + "\"}"
	resp := generalLambda(arnFuncDeleteP2P, paramClientP2P)
	if resp == okC {
		ping := getPingP2P()
		if ping == okC && needCheck == true {
			deleteTargetP2P(false)
		}
	} else if resp == koC {
		log.Fatal(errorDeleteP2P)
	}
}

/*
Read rw pointer to check blockchain
*/
func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil || str == "" {
			return
		} else if str != "\n" {
			chain := make([]Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				fmt.Printf(exceptionJSON)
				log.Fatal(err)
			}
			mutex.Lock()
			//Updating blc from broadcast
			Blockchain = updateBlc(chain, Blockchain)
			mutex.Unlock()
		}
	}
}

/*
Func to write in rw
*/
func writeData(rw *bufio.ReadWriter) {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			closeCon()
			log.Fatalln(sig)
		}
	}()

	//updates Blockchain every 5 seconds in rw
	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				fmt.Printf(exceptionJSON)
				fmt.Println(err)
			}
			mutex.Unlock()
			//spew.Dump(Blockchain)
			mutex.Lock()
			//sending blockchain to broadcast
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			logEntry(blockchainstr, "", 3)
			mutex.Unlock()

		}
	}()

	//go to first menu
	for {
		inOk := false
		for inOk == false {
			inOk = mainActions(rw)
		}
	}
}
