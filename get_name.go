package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

func generateRandomString(length int) string {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Define the characters that can appear in the string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a slice to hold the random string
	result := make([]byte, length)

	// Loop to select random characters
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}

	// Convert the slice to a string and return
	return string(result)
}

// func GetEncryptedFilePath(filePath string, videoUrl string) string {
// 	id := (strings.Split(videoUrl[strings.Index(videoUrl, "v=")+2:len(videoUrl)-1], ""))

// 	for i, j := 0, len(id)-1; i < j; i, j = i+1, j-1 {
// 		id[i], id[j] = id[j], id[i]
// 	}

// 	videoName := strings.Join(id, "")

// 	finalPath := strings.Replace(filePath, filePath[strings.LastIndex(filePath, "\\"):], fmt.Sprintf("\\%s.enc", videoName), 1)

// 	return finalPath
// }

func GetEncryptedFilePath(filePath string, videoName string) string {
	var name string
	if videoName == "" {
		name = generateRandomString(10)
	} else {
		name = videoName
	}

	finalPath := strings.Replace(filePath, filePath[strings.LastIndex(filePath, "\\"):], fmt.Sprintf("\\%s", name), 1)

	return finalPath
}
