package encryptor

import (
	"archive/zip"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var ignoredFiles = []string{
	"manifest.json",
	"pack_icon.png",
}

func EncryptPack(path string, dest string, key string) {
	if _, err := os.Stat(path); err != nil {
		panic(err)
	}
	oldPack, err := zip.OpenReader(path)
	if err != nil {
		panic(err)
	}

	newPackFile, err := os.Create(dest)
	newPack := zip.NewWriter(newPackFile)
	newPack.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	var contents Contents

	for _, f := range oldPack.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if isIgnored(f.Name) {
			ignoredFile, err := newPack.Create(f.Name)
			if err != nil {
				panic(err)
			}
			ignoredReader, err := oldPack.Open(f.Name)
			if err != nil {
				panic(err)
			}
			ignoredBytes, err := io.ReadAll(ignoredReader)
			if err != nil {
				panic(err)
			}
			_, err = ignoredFile.Write(ignoredBytes)
			if err != nil {
				panic(err)
			}
			continue
		}
		fw, err := newPack.Create(f.Name)
		if err != nil {
			panic(err)
		}
		file, err := oldPack.Open(f.Name)
		if err != nil {
			panic(err)
		}
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		fileKey := randomString(32)
		contents.Content = append(contents.Content, Content{Path: f.Name, Key: fileKey})
		encryptedBytes := encrypt(fileBytes, []byte(fileKey))
		fmt.Printf("Encrypting %s with key %s...\n", f.Name, fileKey)
		_, err = fw.Write(encryptedBytes)
		if err != nil {
			panic(err)
		}
	}
	contentJson, err := json.Marshal(&contents)
	if err != nil {
		panic(err)
	}
	contentJsonEncrypted := encrypt(contentJson, []byte(key))
	content, err := newPack.Create("contents.json")
	if err != nil {
		panic(err)
	}
	_, err = content.Write(contentJsonEncrypted)
	if err != nil {
		panic(err)
	}
	err = oldPack.Close()
	if err != nil {
		panic(err)
	}
	err = newPack.Close()
	if err != nil {
		panic(err)
	}
	err = newPackFile.Close()
	if err != nil {
		panic(err)
	}
}

func isIgnored(x string) bool {
	for _, y := range ignoredFiles {
		if x == y {
			return true
		}
	}
	return false
}

func encrypt(content []byte, key []byte) []byte {
	iv := key[:16]
	cipherCtx, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(cipherCtx, iv)
	stream.XORKeyStream(content, content)
	return content
}
