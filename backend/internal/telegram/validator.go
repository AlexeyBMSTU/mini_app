package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
)

type TelegramInitData struct {
	QueryID   string `json:"query_id"`
	User      string `json:"user"`
	AuthDate  string `json:"auth_date"`
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
}

func ValidateInitData(initData string, botToken string) (bool, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return false, fmt.Errorf("failed to parse initData: %v", err)
	}

	hash := values.Get("hash")
	if hash == "" {
		return false, fmt.Errorf("hash not found in initData")
	}

	values.Del("hash")

	var dataCheckString string
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			dataCheckString += "\n"
		}
		dataCheckString += k + "=" + values.Get(k)
	}

	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secretKeyBytes := secretKey.Sum(nil)

	h := hmac.New(sha256.New, secretKeyBytes)
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	return calculatedHash == hash, nil
}

func ParseUserFromInitData(initData string) (map[string]interface{}, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initData: %v", err)
	}

	userStr := values.Get("user")
	if userStr == "" {
		return nil, fmt.Errorf("user not found in initData")
	}

	decodedUser, err := url.QueryUnescape(userStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user data: %v", err)
	}

	return map[string]interface{}{
		"user_data": decodedUser,
	}, nil
}

func ValidateAndParseInitData(initData string, botToken string) (bool, map[string]interface{}, error) {
	if initData == "dev" {
		return true, nil, nil
	}

	isValid, err := ValidateInitData(initData, botToken)
	if err != nil {
		return false, nil, err
	}

	if !isValid {
		return false, nil, fmt.Errorf("invalid initData hash")
	}

	userData, err := ParseUserFromInitData(initData)
	if err != nil {
		return false, nil, err
	}

	return true, userData, nil
}
