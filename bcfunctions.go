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
	record := string(block.Index) + block.Timestamp + string(block.CustomValue) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
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
func generateBlock(oldBlock Block, CustomValue int) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.CustomValue = CustomValue
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

func updateBlc(chain []Block, Blockchain []Block) []Block {
	if len(chain) > len(Blockchain) {
		Blockchain = chain
		//fmt.Printf(greenColorCMD, string(bytes))
	}
	return Blockchain
}

func updateGlobal(bytes []byte) {
	bytestofile(bytes)
	uploadfile()
}

func insertBlc(customvalue int, Blockchain []Block) []Block {
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], customvalue)
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		mutex.Lock()
		Blockchain = append(Blockchain, newBlock)
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
