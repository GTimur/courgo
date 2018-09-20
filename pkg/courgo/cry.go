/*
  Реализует шифрование строковых данных
*/
package courgo

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
	"fmt"
)

/* Вектор инициализации (16 byte) */
var commonIV = []byte("KRASNODAR201703.")
/* Ключ шифрования для AES (32 byte) */
var word = "storeConfigInSafePlace!@#$!@#!@#"

func AesEncrypt(text string) ([]byte, error) {
	key := word
	IV := commonIV
	// Create the aes encryption algorithm
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Printf("Error: NewCipher(%d bytes) = %s\n", len(key), err)
		return nil, err
	}

	// Encrypted string
	cfb := cipher.NewCFBEncrypter(c, IV)
	ciphertext := make([]byte, len(text))
	cfb.XORKeyStream(ciphertext, []byte(text))
	return ciphertext, nil
}

func AesDecript(text []byte) (string, error) {
	key := word
	IV := commonIV
	// Create the aes encryption algorithm
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Printf("Error: NewCipher(%d bytes) = %s\n", len(key), err)
		return "", err
	}

	// Decrypt strings
	cfbdec := cipher.NewCFBDecrypter(c, IV)
	plaintextCopy := make([]byte, len(text))
	cfbdec.XORKeyStream(plaintextCopy, []byte(text))
	return fmt.Sprintf("%s", string(plaintextCopy)), nil
}
