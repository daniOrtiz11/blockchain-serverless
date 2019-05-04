package main

import (
	"crypto/sha256"
	"encoding/hex"
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
	for i := 0; i < len(Bank); i++ {
		if Bank[i].Name == name {
			iacc = i
		}
	}
	return iacc
}

func searchAccountByPublicKey(key string) int {
	iacc := -1
	for i := 0; i < len(Bank); i++ {
		if Bank[i].PublicID == key {
			iacc = i
		}
	}
	return iacc
}

func searchAccountByPrivKey(privateKey string) int {
	var iacc int
	iacc = -1
	for i := 0; i < len(Bank); i++ {
		if Bank[i].PrivateID == privateKey {
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
	newAccount.Amount = initAmountAccount
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
		//prepareUpload(1)
	}
	return Blockchain
}

func updateBank(newMovs int, chain []Block) {
	i := 0
	restartAmountBank()
	for i = 0; i < newMovs; i++ {
		t := chain[i].Transaction
		if t.Amount > -2 {
			indexSource := searchAccountByPublicKey(t.SourceID)
			if indexSource == -1 {
				var acc Account
				acc.PublicID = t.SourceID
				acc.Amount = initAmountAccount
				Bank = append(Bank, acc)
				indexSource = searchAccountByPublicKey(t.SourceID)
			}
			if t.Amount == createAmountName {
				Bank[indexSource].Name = t.TargetID
			} else if t.Amount == createAmountPriv {
				Bank[indexSource].PrivateID = t.TargetID
			} else {
				indexTarget := searchAccountByPublicKey(t.TargetID)
				Bank[indexSource].Amount = Bank[indexSource].Amount - t.Amount
				Bank[indexTarget].Amount = Bank[indexTarget].Amount + t.Amount
				if account.PublicID == Bank[indexTarget].PublicID {
					account.Amount = account.Amount + t.Amount
				} else if account.PublicID == Bank[indexSource].PublicID {
					account.Amount = account.Amount - t.Amount
				}
			}
		}
	}
}

func updateGlobal(bytestoUpload []byte, localfile string, bucketfile string) {
	bytestofile(bytestoUpload, localfile)
	uploadfile(localfile, bucketfile)
}

func insertBlc(transaction Transaction, Blockchain []Block) []Block {
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], transaction)
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		mutex.Lock()
		Blockchain = append(Blockchain, newBlock)
		updateBankByTransaction(transaction)
		prepareUpload(0)
		mutex.Unlock()
	}
	return Blockchain
}

func isUserInBank(id string) bool {
	ok := false
	for _, u := range Bank {
		if u.PublicID == id {
			ok = true
		}
	}
	return ok
}

func getUserByID(id string) Account {
	var us Account
	for _, u := range Bank {
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
		for _, ac := range Bank {
			if ac.PublicID == t.SourceID {
				iS = index
			} else if ac.PublicID == t.TargetID {
				iT = index
			}
			index++
		}
		Bank[iS].Amount = Bank[iS].Amount - t.Amount
		Bank[iT].Amount = Bank[iT].Amount + t.Amount
		prepareUpload(1)
	}
}
