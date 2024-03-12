package hellsgopher

import (
	"crypto/aes"
	"crypto/cipher"
	cr "crypto/rand"
	"encoding/base64"
	"errors"
)

func DecryptString(key []byte, secure string) (decoded string, err error) {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		return
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}

// generate a secure []byte of length l
func GenerateSecureBytes(l int) ([]byte, error) {
	randBytes := make([]byte, l)

	_, err := cr.Read(randBytes)
	if err != nil {
		return nil, err
	}

	return randBytes, nil
}

// generate a 32 byte secure key
func GenerateKey() ([]byte, error) {
	return GenerateSecureBytes(32)
}
