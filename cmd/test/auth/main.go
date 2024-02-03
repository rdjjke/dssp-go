package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func main() {
	sharedSecretKey := make([]byte, 32)
	hash := hmac.New(sha256.New, sharedSecretKey)

	aliceNonce := make([]byte, 24)
	bobNonce := make([]byte, 24)
	aliceID := make([]byte, 8)
	bobID := make([]byte, 8)

	var data []byte
	data = append(data, aliceNonce...)
	data = append(data, bobNonce...)
	data = append(data, aliceID...)
	data = append(data, bobID...)

	hash.Write(data)

	dataHMAC := hash.Sum(nil)

	fmt.Println(len(dataHMAC))
}
