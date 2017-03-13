package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"

)

var commonIV = []byte("KRASNODAR201703.")

func main() {
	// Need to encrypt a string
	plaintext := []byte("PASSWORD")
	// If there is an incoming string of words to be encrypted, set plaintext to that incoming string
	if len(os.Args) > 1 {
		plaintext = []byte(os.Args[1])
	}

	// aes encryption string
	key_text := "storeConfigInSafePlace!@#$!@#!@#"
	if len(os.Args) > 2 {
		key_text = os.Args[2]
	}

	fmt.Println(len(key_text))

	// Create the aes encryption algorithm
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
		os.Exit(-1)
	}

	// Encrypted string
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	fmt.Printf("%s=>%x\n", plaintext, ciphertext)
	fmt.Printf("%d\n",ciphertext)


	// Decrypt strings
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(plaintext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	fmt.Printf("%x=>%s\n", ciphertext, plaintextCopy)
}
