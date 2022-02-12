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
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
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

	inPath := os.Args[1]
	inBytes, err := os.ReadFile(inPath)
	if err != nil {
		log.Fatalf("cannot read input file: %v\n", err)
	}

	outPath := os.Args[2]
	isUpload := false
	if strings.HasPrefix(strings.ToLower(outPath), "https://") {
		log.Printf("Uploading result to %s.\n", outPath)
		isUpload = true
	}

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

	if isUpload {
		u, err := url.Parse(outPath)
		if err != nil {
			log.Fatalf("cannot parse url: %s\n", err)
		}

		if strings.HasSuffix(u.Path, "/") {
			_, filename := filepath.Split(inPath)
			u.Path = path.Join(u.Path, filename)
		}

		req, err := http.NewRequest(http.MethodPut, u.String(), strings.NewReader(finalDocument))
		if err != nil {
			log.Fatalf("cannot build request: %v\n", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("cannot upload file: %v\n", err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			log.Fatalf("unexpected status during upload: %d\n", resp.StatusCode)
		}

		u.User = nil
		u.Fragment = hex.EncodeToString(key)

		fmt.Printf("URL: %s\n", u)
	} else {
		if err := os.WriteFile(outPath, []byte(finalDocument), 0644); err != nil {
			log.Fatalf("cannot write to file: %v\n", err)
		}

		log.Printf("output written to file %s\n", outPath)
		fmt.Printf("Key: %s\n", hex.EncodeToString(key))
	}
}
