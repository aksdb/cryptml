package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "embed"
)

//go:embed webstub/cryptml.html
var stub string

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <input file> <output file>\n", os.Args[0])
		os.Exit(1)
	}

	inBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("cannot read input file: %v\n", err)
	}
	outFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalf("cannot open output file: %v\n", err)
	}
	defer outFile.Close()

	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Fatalf("cannot initialize key: %v\n", err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("cannot build cipher: %v\n", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("cannot build cipher: %v\n", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("cannot initialize nonce: %v\n", err)
	}

	b := gcm.Seal(nonce, nonce, inBytes, nil)

	finalDocument := strings.Replace(stub, "{{PLACEHOLDER}}", base64.StdEncoding.EncodeToString(b), 1)
	outFile.Write([]byte(finalDocument))

	fmt.Printf("Key: %s\n", hex.EncodeToString(key))
}
