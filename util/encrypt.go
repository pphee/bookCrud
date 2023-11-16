package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type IEncryptionService interface {
	Encrypt(text string) (string, error)
	Decrypt(cryptoText string) (string, error)
}

type EncryptionService struct {
	key []byte
}

func NewEncryptionService(key []byte) IEncryptionService {
	return &EncryptionService{key: key}
}

func (s *EncryptionService) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(s.key)

	if err != nil {
		return "", err
	}

	paddedText := pkcs7Padding([]byte(text), aes.BlockSize)
	ciphertext := make([]byte, len(paddedText))

	mode := newECEncrypted(block)
	mode.CryptBlocks(ciphertext, paddedText)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *EncryptionService) Decrypt(cryptoText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := newECDecrypted(block)
	mode.CryptBlocks(plaintext, ciphertext)

	unpaidPlaintext, err := pkcs7Unpadding(plaintext)
	if err != nil {
		return "", err
	}

	return string(unpaidPlaintext), nil
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	latest := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, latest...)
}

func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: Data is empty")
	}

	padLen := int(data[length-1])
	if padLen > length || padLen > aes.BlockSize {
		return nil, errors.New("pkcs7: Invalid padding")
	}

	for _, padByte := range data[length-padLen:] {
		if int(padByte) != padLen {
			return nil, errors.New("pkcs7: Invalid padding")
		}
	}

	return data[:length-padLen], nil
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecEncrypted ecb

func newECEncrypted(b cipher.Block) cipher.BlockMode {
	return (*ecEncrypted)(newECB(b))
}

func (x *ecEncrypted) BlockSize() int { return x.blockSize }

func (x *ecEncrypted) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecDecrypted ecb

func newECDecrypted(b cipher.Block) cipher.BlockMode {
	return (*ecDecrypted)(newECB(b))
}

func (x *ecDecrypted) BlockSize() int { return x.blockSize }

func (x *ecDecrypted) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
