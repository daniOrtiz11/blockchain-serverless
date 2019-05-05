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

/*
Func to view Blockchain
*/
func viewState(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == okC {
		fmt.Println(options1Title1)
		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			fmt.Println(err)
		}
		spew.Dump(Blockchain)
		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
	} else {
		fmt.Println(koPingP2P)
		fmt.Println(reconnectingP2P)
		generalMain()
	}
}

/*
Func to view Bank
*/
func viewBank(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == okC {
		fmt.Println(options4Ttile1)
		mutex.Lock()
		for i := 0; i < len(Bank); i++ {
			acc := Bank[i]
			index := i + 1
			fmt.Println(strconv.Itoa(index) + ". " + acc.Name + " " + strconv.Itoa(acc.Amount) + " " + string(acc.PublicID))
		}
		mutex.Unlock()
	} else {
		fmt.Println(koPingP2P)
		fmt.Println(reconnectingP2P)
		generalMain()
	}
}

/*
Func to view your account
*/
func viewAccountState(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == okC {
		fmt.Println(options2Title1)
		mutex.Lock()
		toStringAccount()
		mutex.Unlock()
	} else {
		fmt.Println(koPingP2P)
		fmt.Println(reconnectingP2P)
		generalMain()
	}
}

/*
Func to end session
*/
func closeCon() {
	deleteTargetP2P(true)
	log.Fatal(endMessage)
}

/*
Func to do a transfer
*/
func insertBlock() {
	pingP2P := getPingP2P()
	if pingP2P == okC {
		enabledAmount := false
		realAmount := 0
		for enabledAmount == false {
			fmt.Println(options3Title11)
			fmt.Print("> ")
			stdReader := bufio.NewReader(os.Stdin)
			sendData, err := stdReader.ReadString('\n')
			if err != nil {
				fmt.Printf(exceptionReader)
				log.Fatal(err)
			}
			sendData = strings.Replace(sendData, "\n", "", -1)
			//check if is a number
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
				fmt.Printf(exceptionReader)
				log.Fatal(err)
			}
			idUser := strings.Replace(sendData, "\n", "", -1)
			if err == nil {
				//check if target user exists
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
		//do transfer
		var t Transaction
		t.Amount = realAmount
		t.SourceID = account.PublicID
		t.TargetID = userTarget.PublicID
		account.Amount = account.Amount - t.Amount
		Blockchain = insertBlc(t, Blockchain)
		//send to upload
		prepareUpload(0)
	} else {
		fmt.Println(koPingP2P)
		fmt.Println(reconnectingP2P)
		generalMain()
	}
}

/*
Func to insert some blocks to create a new account
*/
func insertAccountBlock(newac Account) {
	var t Transaction
	t.Amount = createAmountName
	t.SourceID = newac.PublicID
	t.TargetID = newac.Name
	Blockchain = insertBlc(t, Blockchain)
	t.Amount = createAmountPriv
	t.SourceID = newac.PublicID
	t.TargetID = newac.PrivateID
	Blockchain = insertBlc(t, Blockchain)
	prepareUpload(0)
}

/*
Func to login
*/
func login(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == okC {
		fmt.Println(login1Title1)
		fmt.Print("> ")
		stdReader := bufio.NewReader(os.Stdin)
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Printf(exceptionReader)
			log.Fatal(err)
		}
		privateKey := strings.Replace(sendData, "\n", "", -1)
		if err == nil {
			//check if privateKey exists
			iacc := searchAccountByPrivKey(privateKey)
			if iacc == -1 {
				fmt.Println(errorLogin)
			} else {
				//set your account
				account = Bank[iacc]
				logged = true
				fmt.Println(welcomeLogged + account.Name + "!")
				//go to nexts actions
				loggedActions(rw, stdReader)
			}
		}
	} else {
		fmt.Println(koPingP2P)
		fmt.Println(reconnectingP2P)
		generalMain()
	}
}

/*
Func to create an Account
*/
func createAccount(rw *bufio.ReadWriter) {
	pingP2P := getPingP2P()
	if pingP2P == okC {
		usedName := -1
		for usedName == -1 {
			fmt.Println(login2Title1)
			fmt.Print("> ")
			stdReader := bufio.NewReader(os.Stdin)
			sendData, err := stdReader.ReadString('\n')
			if err != nil {
				fmt.Printf(exceptionReader)
				log.Fatal(err)
			}
			sendData = strings.Replace(sendData, "\n", "", -1)
			//check if name exists in Bank
			usedName = searchAccountByName(sendData)
			if usedName != -1 {
				fmt.Println(errorCreateAccount)
			} else {
				Bank = insertAccount(sendData, Bank)
				//go to upload
				prepareUpload(1)
				fmt.Println(login2Title2)
				fmt.Println(login2Title3)
				toStringAccount()
				logged = true
				loggedActions(rw, stdReader)
			}
		}
	} else {
		fmt.Println(koPingP2P)
		fmt.Println(reconnectingP2P)
		generalMain()
	}
}

/*
Func to show logged actions. Only allow if user is logged
*/
func loggedActions(rw *bufio.ReadWriter, stdReader *bufio.Reader) {
	for {
		inOk := false
		for inOk == false {
			showMenu2()
			sendData, err := stdReader.ReadString('\n')
			if err != nil {
				fmt.Printf(exceptionReader)
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

/*
Func to show first actions. Only allow if user is not logged
*/
func mainActions(rw *bufio.ReadWriter) bool {
	inOk := false
	stdReader := bufio.NewReader(os.Stdin)
	if logged == false {
		showMenu1()
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Printf(exceptionReader)
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
	} else {
		loggedActions(rw, stdReader)
	}

	return inOk
}

func restartAmountBank() {
	for i := 0; i < len(Bank); i++ {
		Bank[i].Amount = initAmountAccount
	}
	account.Amount = initAmountAccount
}
