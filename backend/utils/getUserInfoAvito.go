package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"io"
)

type UserInfoAvitoResponse struct {
		Email 			string `json:"email"`
		Id   				int32 `json:"id"`
		Name   			string    `json:"name"`
		Phone   		string    `json:"phone"`
		Phones   		[]string    `json:"phones"`
		ProfileUrl   string    `json:"profile_url"`
}

func GetUserInfoAvito() (UserInfoAvitoResponse, error) {
	reqURL := "https://api.avito.ru/core/v1/accounts/self"
	
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Printf("Error creating request to external API: %v", err)
		return UserInfoAvitoResponse{}, err
	}

	token, _ := GetToken()
	log.Printf("Token: %s", string(token.AccessToken))

	authHeader := fmt.Sprintf("%s %s", token.TokenType, token.AccessToken)

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to external API: %v", err)
		return UserInfoAvitoResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from external API: %v", err)
		return UserInfoAvitoResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("External API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
		return UserInfoAvitoResponse{}, err
	}

	defer resp.Body.Close();

	var userInfoAvitoResponse struct {
			Email 			string `json:"email"`
			Id   				int32 `json:"id"`
			Name   			string    `json:"name"`
			Phone   		string    `json:"phone"`
			Phones   		[]string    `json:"phones"`
			ProfileUrl   string    `json:"profile_url"`
	}

	if err := json.Unmarshal(body, &userInfoAvitoResponse); err != nil {
			return UserInfoAvitoResponse{}, fmt.Errorf("error unmarshaling token response: %v", err)
	}

	return userInfoAvitoResponse, nil
}
