package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func EncryptFile(filePath string, encryptedPath string, secretKey string, iv string) error {
	// Convert the secret key and IV from hex strings to byte slices
	key, err := hex.DecodeString(secretKey)
	if err != nil {
		return fmt.Errorf("failed to decode secret key: %v", err)
	}

	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return fmt.Errorf("failed to decode IV: %v", err)
	}

	// Create a new AES cipher block and CBC encrypter
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %v", err)
	}

	cipherStream := cipher.NewCBCEncrypter(block, ivBytes)

	// Open the input file
	inputFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	// Open the output file
	outputFile, err := os.Create(encryptedPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Create buffers for reading from the input file and writing to the output
	buffer := make([]byte, aes.BlockSize*1024)
	for {
		bytesRead, readErr := inputFile.Read(buffer)
		if bytesRead > 0 {
			data := buffer[:bytesRead]

			// Apply padding if it's the last read
			if bytesRead < len(buffer) || readErr == io.EOF {
				data = pkcs7Padding(data, aes.BlockSize)
			}

			cipherStream.CryptBlocks(data, data)
			if _, writeErr := outputFile.Write(data); writeErr != nil {
				return fmt.Errorf("failed to write to output file: %v", writeErr)
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return fmt.Errorf("failed to read input file: %v", readErr)
		}
	}

	return nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	paddingSize := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(data, padding...)
}
