package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// Okta details
	clientID := "0oamjhi22uOWJauGd697"
	clientSecret := "mrfAcwf0xA41LR4liDBSev-onsMlE5zkdxrECQXKlnILArOPW5KX9w1zCB-8DjNw"
	authURL := "https://trial-2902165-admin.okta.com/oauth2/default/v1/token"

	// Prepare the Authorization header
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret))

	// Prepare the request payload
	data := []byte("grant_type=client_credentials&scope=openid")

	// Create the HTTP request
	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", authHeader)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	if resp.StatusCode == http.StatusOK {
		// Parse the JSON response to extract the access token
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Error parsing JSON response:", err)
			os.Exit(1)
		}

		// Print the access token
		fmt.Println("Access Token:", response["access_token"])
	} else {
		// Print the error response
		fmt.Printf("Failed to get token: %s\nResponse: %s\n", resp.Status, string(body))
	}
}
