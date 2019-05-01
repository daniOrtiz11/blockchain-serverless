package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func bytestofile(b []byte) {
	f, err := os.Create(localfile)
	check(err)
	defer f.Close()
	n2, err := f.Write(b)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parserLocalP2P(fullstr string) {
	splitFull := strings.Split(fullstr, "/")
	localP2P.Ipdir = splitFull[2]
	localP2P.Port = splitFull[4]
	localP2P.Key = splitFull[6]
}

func parserTarget(target string) string {
	target = strings.Replace(target, "\"", "", -1)
	splitFull := strings.Split(target, ":")
	localP2P.PrevKey = splitFull[2]
	localTarget := ipv4_2 + splitFull[0] + tcp2 + splitFull[1] + ipfs2 + splitFull[2]
	return localTarget
}

func toStringTransaction(t Transaction) string {
	return t.SourceID + "##" + strconv.Itoa(t.Amount) + "##" + t.TargetID
}

func toStringAccount() {
	fmt.Println("Name: " + account.Name)
	fmt.Println("Coins: " + strconv.Itoa(account.Amount))
	fmt.Println("Public Key: " + account.PublicID)
	fmt.Println("Private Key: " + account.PrivateID)
}

func showMenu1() {
	fmt.Println(optionsTitle)
	fmt.Println(login1)
	fmt.Println(login2)
	fmt.Println(closeOptionLogin)
	fmt.Print("> ")
}

func showMenu2() {
	fmt.Println(optionsTitle)
	fmt.Println(options1Title)
	fmt.Println(options2Title)
	fmt.Println(options3Title)
	fmt.Println(options4Title)
	fmt.Println(closeOptionMenu)
	fmt.Print("> ")
}
