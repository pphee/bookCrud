package util

//
//import (
//	"crypto/aes"
//	"crypto/cipher"
//	"crypto/rand"
//	"encoding/base64"
//	"errors"
//	"fmt"
//	"io"
//)
//
//type IEncryptionService interface {
//	Encrypt(text string) (string, error)
//	Decrypt(cryptoText string) (string, error)
//}
//
//type EncryptionService struct {
//	key []byte
//}
//
//func NewEncryptionService(key []byte) IEncryptionService {
//	return &EncryptionService{key: key}
//}
//
//func (s *EncryptionService) Encrypt(text string) (string, error) {
//	block, err := createCipherBlock(s.key)
//	if err != nil {
//		return "", err
//	}
//	b := base64.StdEncoding.EncodeToString([]byte(text))
//	ciphertext := make([]byte, aes.BlockSize+len(b))
//	iv := ciphertext[:aes.BlockSize]
//	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
//		return "", err
//	}
//	stream := cipher.NewCFBEncrypter(block, iv)
//	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
//	return base64.StdEncoding.EncodeToString(ciphertext), nil
//}
//
//func (s *EncryptionService) Decrypt(cryptoText string) (string, error) {
//	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
//	if err != nil {
//		return "", err
//	}
//	block, err := createCipherBlock(s.key)
//	if err != nil {
//		return "", err
//	}
//	if len(ciphertext) < aes.BlockSize {
//		return "", errors.New("ciphertext too short")
//	}
//	iv := ciphertext[:aes.BlockSize]
//	ciphertext = ciphertext[aes.BlockSize:]
//	stream := cipher.NewCFBDecrypter(block, iv)
//	stream.XORKeyStream(ciphertext, ciphertext)
//	data, err := base64.StdEncoding.DecodeString(string(ciphertext))
//	if err != nil {
//		return "", err
//	}
//	return string(data), nil
//}
//
//func createCipherBlock(key []byte) (cipher.Block, error) {
//	fmt.Println(len(key))
//	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
//		return nil, errors.New("invalid key size")
//	}
//	return aes.NewCipher(key)
//}
