package main

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
)

func main() {
	msg := []byte("Hello, world!")
	log.Printf("Source message: %q", msg)

	key := make([]byte, 16)
	for i := range key {
		key[i] = byte(i)
	}

	encrypted := encrypt(msg, key)
	log.Printf("Encrypted message: %v", encrypted)

	decrypted := decrypt(encrypted, key)
	log.Printf("Decrypted message: %q", decrypted)
}

func encrypt(msg, key []byte) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("New AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("New GCM: %v", err)
	}

	log.Printf("Nonce size: %d", gcm.NonceSize())

	nonce := make([]byte, gcm.NonceSize())

	return gcm.Seal(nonce, nonce, msg, nil)
}

func decrypt(msg, key []byte) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("New AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("New GCM: %v", err)
	}

	nonce, msg := msg[:gcm.NonceSize()], msg[gcm.NonceSize():]

	res, err := gcm.Open(nil, nonce, msg, nil)
	if err != nil {
		log.Fatalf("Decrypt message: %v", err)
	}
	return res
}
