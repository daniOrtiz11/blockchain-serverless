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

func handleStream(s net.Stream) {

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

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
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf(ipfs, basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
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
	/*
		Handle clase if there are any problem get addr:
	*/
	if findaddr == false {
		addrcero := basicHost.Addrs()[0]
		fullAddrcero := addrcero.Encapsulate(hostAddr)
		sfulladdr := fmt.Sprintf("%s", fullAddrcero)
		s := strings.Split(sfulladdr, "/")
		nextkey := s[3]
		if nextkey == tcp {
			nextkey = s[6]
		}
		addrstr := fmt.Sprintf(iplocalhost, listenPort)
		realaddress := addrstr + ipfs2 + nextkey
		parserLocalP2P(realaddress)
	} else {
		fullAddr := addr.Encapsulate(hostAddr)
		parserLocalP2P(fullAddr.String())
	}

	return basicHost, nil
}

func getTargetP2P() string {
	target := generalLambda(arnFuncGetP2P, "")
	targetParser := ""
	if target != "" && target != "empty" {
		targetParser = parserTarget(target)
	}
	return targetParser
}

func getPingP2P() string {
	clientP2P := localP2P.Ipdir + "/" + localP2P.Port + "/" + localP2P.Key + "/" + localP2P.PrevKey
	paramClientP2P := "{\"newnode\": \"" + clientP2P + "\"}"
	ping := generalLambda(arnFuncPingP2P, paramClientP2P)
	pingResult := "ko"
	if ping == "ok" {
		pingResult = "ok"
	}
	return pingResult
}

func setTargetP2P() {
	//ip:port:key
	clientP2P := localP2P.Ipdir + "/" + localP2P.Port + "/" + localP2P.Key + "/" + localP2P.PrevKey
	paramClientP2P := "{\"newnode\": \"" + clientP2P + "\"}"
	resp := generalLambda(arnFuncSetP2P, paramClientP2P)
	if resp == "ok" {
		setTargetP2P()
	} else if resp == "ko" {
		log.Fatal(errorSetP2P)
	}
}

func deleteTargetP2P(needCheck bool) {
	//ip:port:key
	clientP2P := localP2P.Ipdir + "/" + localP2P.Port + "/" + localP2P.Key + "/" + localP2P.PrevKey
	paramClientP2P := "{\"newnode\": \"" + clientP2P + "\"}"
	resp := generalLambda(arnFuncDeleteP2P, paramClientP2P)
	if resp == "ok" {
		ping := getPingP2P()
		if ping == "ok" && needCheck == true {
			deleteTargetP2P(false)
		}
	} else if resp == "ko" {
		log.Fatal(errorDeleteP2P)
	}
	//ok or ko
}

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
			//spew.Dump(Blockchain)
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
			inOk = mainActions(rw, stdReader)
		}
	}
}
