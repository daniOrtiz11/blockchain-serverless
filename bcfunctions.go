package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"
)

// SHA256 hasing
func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.Transaction.SourceID+block.Transaction.TargetID) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// SHA256 hasing
func calculateHashAccount(seed string) string {
	record := string(seed)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	hashstr := hex.EncodeToString(hashed)
	return hashstr
}

func searchAccountByName(name string) int {
	iacc := -1
	for i := 0; i < len(bank); i++ {
		if bank[i].Name == name {
			iacc = i
		}
	}
	return iacc
}

func searchAccountByPublicKey(key string) int {
	iacc := -1
	for i := 0; i < len(bank); i++ {
		if bank[i].PublicID == key {
			iacc = i
		}
	}
	return iacc
}

func searchAccountByPrivKey(privateKey string) int {
	var iacc int
	iacc = -1
	for i := 0; i < len(bank); i++ {
		if bank[i].PrivateID == privateKey {
			iacc = i
		}
	}
	return iacc
}

func genesisTransaction() Transaction {
	transaction := Transaction{}

	return transaction
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, transaction Transaction) Block {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transaction = transaction
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

// create a new block using previous block's hash
func generateAccount(name string) Account {
	var newAccount Account
	newAccount.Name = name
	newAccount.Amount = 100
	newAccount.PublicID = calculateHashAccount(name)
	t := time.Now()
	newAccount.PrivateID = calculateHashAccount(name + string(newAccount.PublicID) + t.String())
	return newAccount
}

func updateBlc(chain []Block, Blockchain []Block) []Block {
	if len(chain) > len(Blockchain) {
		diflen := len(chain)
		Blockchain = chain
		updateBank(diflen, chain)
	}
	return Blockchain
}

func updateBank(newMovs int, chain []Block) {
	i := 0
	for i = 0; i < newMovs; i++ {
		t := chain[i].Transaction
		if t.Amount > -2 {
			indexSource := searchAccountByPublicKey(t.SourceID)
			if indexSource == -1 {
				var acc Account
				acc.PublicID = t.SourceID
				acc.Amount = 100
				bank = append(bank, acc)
				indexSource = searchAccountByPublicKey(t.SourceID)
			}
			if t.Amount == -1 {
				bank[indexSource].Name = t.TargetID
			} else if t.Amount == 0 {
				bank[indexSource].PrivateID = t.TargetID
			} else {
				indexTarget := searchAccountByPublicKey(t.TargetID)
				bank[indexSource].Amount = bank[indexSource].Amount - t.Amount
				bank[indexTarget].Amount = bank[indexTarget].Amount + t.Amount
			}
		}
	}
}

func updateGlobal(bytes []byte) {
	bytestofile(bytes)
	uploadfile()
}

func insertBlc(transaction Transaction, Blockchain []Block) []Block {
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], transaction)
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		mutex.Lock()
		Blockchain = append(Blockchain, newBlock)
		updateBankByTransaction(transaction)
		bytes, err := json.MarshalIndent(Blockchain, "", "  ")
		if err != nil {
			log.Printf(exceptionJSON)
			log.Fatal(err)
		}
		debug := false
		if debug == true {
			updateGlobal(bytes)
		}
		mutex.Unlock()
	}
	return Blockchain
}

func isUserInBank(id string) bool {
	ok := false
	for _, u := range bank {
		if u.PublicID == id {
			ok = true
		}
	}
	return ok
}

func getUserByID(id string) Account {
	var us Account
	for _, u := range bank {
		if u.PublicID == id {
			us = u
		}
	}
	return us
}

func insertAccount(name string, bank []Account) []Account {
	newAccount := generateAccount(name)
	insertAccountBlock(newAccount)
	bank = append(bank, newAccount)
	account = newAccount
	return bank
}

func updateBankByTransaction(t Transaction) {
	iS := 0
	iT := 0
	index := 0
	if t.Amount > 0 {
		for _, ac := range bank {
			if ac.PublicID == t.SourceID {
				iS = index
			} else if ac.PublicID == t.TargetID {
				iT = index
			}
			index++
		}
		bank[iS].Amount = bank[iS].Amount - t.Amount
		bank[iT].Amount = bank[iT].Amount + t.Amount
	}
}
