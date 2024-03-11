package hellsgopher

import (
	"crypto/aes"
	"crypto/cipher"
	cr "crypto/rand"
	"encoding/binary"
)

// decrypt []byte with given key
func DecryptBytes(message []byte, key []byte) ([]byte, error) {
	iv := message[0:16]
	cText := message[16:]

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	c := cipher.NewCBCDecrypter(cipherBlock, iv)
	d := make([]byte, len(cText))
	c.CryptBlocks(d, cText)

	lenBytes := d[0:4]
	len := binary.LittleEndian.Uint32(lenBytes)
	d = d[4:]
	return d[:len], nil
}

// return a decrypted string using given key
func DecryptString(s string, key []byte) (string, error) {
	d, err := DecryptBytes([]byte(s), key)
	if err != nil {
		return "", err
	}

	return string(d), nil
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
