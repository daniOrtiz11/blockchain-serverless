package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

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

func viewState(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == "ok" {
		fmt.Println(options1Title1)
		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			log.Println(err)
		}
		spew.Dump(Blockchain)
		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
	} else {
		log.Println(koPingP2P)
		log.Println(reconnectingP2P)
		generalMain()
	}
}

func viewBank(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == "ok" {
		fmt.Println("Current Bank state:")
		mutex.Lock()
		for i := 0; i < len(bank); i++ {
			acc := bank[i]
			index := i + 1
			fmt.Println(strconv.Itoa(index) + ". " + acc.Name + " " + strconv.Itoa(acc.Amount) + " " + string(acc.PublicID))
		}
		mutex.Unlock()
	} else {
		log.Println(koPingP2P)
		log.Println(reconnectingP2P)
		generalMain()
	}
}

func viewAccountState(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == "ok" {
		fmt.Println("Current Account state:")
		mutex.Lock()
		toStringAccount()
		mutex.Unlock()
	} else {
		log.Println(koPingP2P)
		log.Println(reconnectingP2P)
		generalMain()
	}
}

func closeCon() {
	deleteTargetP2P(true)
	log.Fatal(endMessage)
}

func insertBlock() {
	pingP2P := getPingP2P()
	if pingP2P == "ok" {
		enabledAmount := false
		realAmount := 0
		for enabledAmount == false {
			fmt.Println(options3Title11)
			fmt.Print("> ")
			stdReader := bufio.NewReader(os.Stdin)
			sendData, err := stdReader.ReadString('\n')
			if err != nil {
				log.Printf(exceptionReader)
				log.Fatal(err)
			}
			sendData = strings.Replace(sendData, "\n", "", -1)
			amount, err := strconv.Atoi(sendData)
			if err == nil {
				if amount < account.Amount {
					enabledAmount = true
					realAmount = amount
				} else {
					fmt.Println(options3Title12)
				}
			} else {
				fmt.Println(badFormatNumber)
			}
		}
		isUser := false
		var userTarget Account
		for isUser == false {
			fmt.Println(options3Title21)
			fmt.Print("> ")
			stdReader := bufio.NewReader(os.Stdin)
			sendData, err := stdReader.ReadString('\n')
			if err != nil {
				log.Printf(exceptionReader)
				log.Fatal(err)
			}
			idUser := strings.Replace(sendData, "\n", "", -1)
			if err == nil {
				isUser = isUserInBank(idUser)
				if !isUser {
					fmt.Println(options3Title22)
				} else {
					userTarget = getUserByID(idUser)
				}
			} else {
				fmt.Println(badFormatNumber)
			}
		}
		fmt.Println(options3Title3 + userTarget.Name)
		var t Transaction
		t.Amount = realAmount
		t.SourceID = account.PublicID
		t.TargetID = userTarget.PublicID
		account.Amount = account.Amount - t.Amount
		Blockchain = insertBlc(t, Blockchain)
	} else {
		log.Println(koPingP2P)
		log.Println(reconnectingP2P)
		generalMain()
	}
}

func insertAccountBlock(newac Account) {
	var t Transaction
	t.Amount = -1
	t.SourceID = newac.PublicID
	t.TargetID = newac.Name
	Blockchain = insertBlc(t, Blockchain)
	t.Amount = 0
	t.SourceID = newac.PublicID
	t.TargetID = newac.PrivateID
}

func login(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == "ok" {
		fmt.Println("Introduce your private key:")
		fmt.Print("> ")
		stdReader := bufio.NewReader(os.Stdin)
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Printf(exceptionReader)
			log.Fatal(err)
		}
		privateKey := strings.Replace(sendData, "\n", "", -1)
		if err == nil {
			iacc := searchAccountByPrivKey(privateKey)
			if iacc == -1 {
				fmt.Println("Error in your login. Check your credentials")
			} else {
				account = bank[iacc]
				loggedActions(rw, stdReader)
			}
		}
	} else {
		log.Println(koPingP2P)
		log.Println(reconnectingP2P)
		generalMain()
	}
}

func createAccount(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == "ok" {
		usedName := -1
		for usedName == -1 {
			fmt.Println("Introduce your new name:")
			fmt.Print("> ")
			stdReader := bufio.NewReader(os.Stdin)
			sendData, err := stdReader.ReadString('\n')
			if err != nil {
				log.Printf(exceptionReader)
				log.Fatal(err)
			}
			sendData = strings.Replace(sendData, "\n", "", -1)
			usedName = searchAccountByName(sendData)
			if usedName != -1 {
				fmt.Println("Error. Name used. Choose other.")
			} else {
				bank = insertAccount(sendData, bank)

				fmt.Println("Succesfully creating account!")
				toStringAccount()
				loggedActions(rw, stdReader)
			}
		}
	} else {
		log.Println(koPingP2P)
		log.Println(reconnectingP2P)
		generalMain()
	}
}

func loggedActions(rw *bufio.ReadWriter, stdReader *bufio.Reader) {
	for {
		inOk := false
		for inOk == false {
			showMenu2()
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
					viewAccountState(rw)
					inOk = true
				case 3:
					insertBlock()
					inOk = true
				case 4:
					viewBank(rw)
					inOk = true
				case 5:
					closeCon()
					inOk = true
				default:
					fmt.Println(badFormatOption)
				}
			} else {
				fmt.Println(badFormatNumber)
			}
			fmt.Println()
		}
	}
}

func mainActions(rw *bufio.ReadWriter, stdReader *bufio.Reader) bool {
	inOk := false
	showMenu1()
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
			login(rw)
			inOk = true
		case 2:
			createAccount(rw)
			inOk = true
		case 3:
			closeCon()
			inOk = true
		default:
			fmt.Println(badFormatOption)
		}
	} else {
		fmt.Println(badFormatNumber)
	}
	return inOk
}
