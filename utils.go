package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
Create a file from bytes
*/
func bytestofile(b []byte, localfile string) {
	f, err := os.Create(localfile)
	check(err)
	defer f.Close()
	_, err2 := f.Write(b)
	check(err2)
}

func addtofile(text string, localfile string) {
	f, err := os.OpenFile(localfile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		bytestofile([]byte(text), localfile)
	} else {
		defer f.Close()
		if _, err = f.WriteString(text); err != nil {
			bytestofile([]byte(text), localfile)
		}
	}

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

func toStringAccount() {
	fmt.Println(nameAccountStr + account.Name)
	fmt.Println(coinsAccountStr + strconv.Itoa(account.Amount))
	fmt.Println(publicAccountStr + account.PublicID)
	fmt.Println(privateAccountStr + account.PrivateID)
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

func logEntry(name1 string, name2 string, action int) {
	t := time.Now()
	time := t.Format("20060102150405")
	localfolder := localP2P.Ipdir + ":" + localP2P.Port + "/"
	newEntry := time + logDelimeter + localP2P.Ipdir + ":" + localP2P.Port + logDelimeter
	//action = 0 -> create/write account
	//action = 1 -> get account from blockchain
	//action = 2 -> transaction
	//action = 3 -> create/write in blockchain
	//action = 4 -> udpate blockchain
	switch action {
	case 0:
		newEntry = newEntry + "A" + logDelimeter + localfolder + folder1logs + name1 + logDelimeter + logColor
	case 1:
		newEntry = newEntry + "M" + logDelimeter + localfolder + folder1logs + name1 + logDelimeter + logColor
	case 2:
		newEntry = newEntry + "A" + logDelimeter + localfolder + folder1logs + name1 + logDelimeter + logColor
		newEntry = newEntry + time + logDelimeter + localP2P.Ipdir + ":" + localP2P.Port + logDelimeter +
			"A" + logDelimeter + localfolder + folder1logs + name2 + logDelimeter + logColor
	case 3:
		newEntry = newEntry + "A" + logDelimeter + name1 + logDelimeter + logColor
	case 4:
		newEntry = newEntry + "M" + logDelimeter + name1 + logDelimeter + logColor
	}
	setLog(newEntry)
}

func readTextFromFile(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	return str
}
