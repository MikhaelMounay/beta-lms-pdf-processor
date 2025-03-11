package main

import (
	"encoding/base64"
	"fmt"

	"resty.dev/v3"
)

// Struct for parsing the v3 API response
type AuthResponse struct {
	AuthToken string `json:"authorizationToken"`
	APIInfo   struct {
		StorageAPI struct {
			APIUrl string `json:"apiUrl"`
		} `json:"storageApi"`
	} `json:"apiInfo"`
}

// GetUploadURL authenticates using Application Key ID and Key, then retrieves an upload URL from Backblaze B2
func GetB2UploadURL(appKeyID, appKey, bucketID string) error {
	client := resty.New()
	
	var authData AuthResponse
	// Step 1: Authenticate and Get Authorization Token
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(appKeyID+":"+appKey))
	resp, err := client.R().
		SetHeader("Authorization", authHeader).
		SetResult(&authData).
		EnableRetryDefaultConditions().
		Get("https://api.backblazeb2.com/b2api/v3/b2_authorize_account")

	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to authorize account: %s", resp.String())
	}

	apiURL := authData.APIInfo.StorageAPI.APIUrl
	authToken := authData.AuthToken

	// Step 2: Get Upload URL
	var uploadData map[string]interface{}
	resp, err = client.R().
		SetHeader("Authorization", authToken).
		SetBody(map[string]string{"bucketId": bucketID}).
		SetResult(&uploadData).
		EnableRetryDefaultConditions().
		Get(apiURL + fmt.Sprintf("/b2api/v3/b2_get_upload_url?bucketId=%s", bucketID))

	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to get upload URL: %s", resp.String())
	}

	UploadURL = uploadData["uploadUrl"].(string)
	AuthToken = uploadData["authorizationToken"].(string)

	return nil
}
