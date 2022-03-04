package goframework

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

const (
	randomString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_+"
)

// randomString generates a random string length n from values in the const randomString
func (g *Goframework) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomString)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	return string(s)
}

func (g *Goframework) createDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Goframework) CreateFileIfNotExists(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}

		// need to pass file to the defered func
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	return nil
}

type Encryption struct {
	Key []byte
}

func (e *Encryption) Encrypt(text string) (string, error) {
	plainText := []byte(text)

	// part of the encryption process
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	// create the dataStructure to use for encryption
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		if err != nil {
			return "", err
		}
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)

	return base64.URLEncoding.EncodeToString(ciphertext), nil

}

func (e *Encryption) Decrypt(cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", err
	}

	// get the last part(of length aes.BlockSize) of the ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// do the decryption
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}
