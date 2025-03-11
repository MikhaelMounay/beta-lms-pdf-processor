package main

import (
	"fmt"
	"io"
	"os"

	"resty.dev/v3"
)

// UploadFile uploads a file to Backblaze B2 (Sequential Upload)
func UploadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	sha1sum, err := ComputeSHA1Hash(filePath)
	if err != nil {
		return fmt.Errorf("failed to compute SHA-1: %v", err)
	}

	// Read entire file into memory (prevents chunked encoding issue)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", AuthToken).
		SetHeader("X-Bz-File-Name", BaseB2Folder+fileInfo.Name()).
		SetHeader("Content-Type", "b2/x-auto").
		SetHeader("X-Bz-Content-Sha1", sha1sum).
		SetHeader("Content-Length", fmt.Sprintf("%d", fileSize)).
		SetBody(fileBytes).
		EnableRetryDefaultConditions().
		Post(UploadURL)

	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("upload failed: %s", resp.String())
	}

	fmt.Printf("Uploaded: %s (%.2f MB)\n", fileInfo.Name(), float64(fileSize)/(1024*1024))

	// Delete file after successful upload
	file.Close()
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("Warning: Failed to delete %s: %v\n", filePath, err)
	}

	return nil
}
