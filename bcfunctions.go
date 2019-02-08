package blc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// SHA256 hasing
func CalculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.Private_value) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	fmt.Println("here")
	return hex.EncodeToString(hashed)
}
