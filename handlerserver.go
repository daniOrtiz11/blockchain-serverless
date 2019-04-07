package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	crypto "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"
	ma "gx/ipfs/QmTZBfrPJmjWsCvHEtX5FE6KimVJhsJg5sBbqEFYf4UZtL/go-multiaddr"
	net "gx/ipfs/QmY3ArotKMKaL7YGfbQfyDrib6RVraLqZYWXZvVgZktBxp/go-libp2p-net"
	host "gx/ipfs/QmYrWiWM4qtrnCeT3R14jY3ZZyirDNJgwK57q4qFYePgbd/go-libp2p-host"
	"io"
	mrand "math/rand"
	"strings"

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
		if debug {
			cmdConsole = fmt.Sprintf(debugcmd, listSources, listenPort+1, realaddress)
		} else {
			cmdConsole = fmt.Sprintf(prodcmd, listenPort+1, realaddress)
		}
	} else {
		fullAddr := addr.Encapsulate(hostAddr)
		if debug {
			cmdConsole = fmt.Sprintf(debugcmd, listSources, listenPort+1, fullAddr)
		} else {
			cmdConsole = fmt.Sprintf(prodcmd, listenPort+1, fullAddr)
		}
	}

	return basicHost, nil
}
