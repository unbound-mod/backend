package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	Structures "unbound.rip/backend/structures"
)

func GetAuthorizationFromCode(code string, redirect string) (*Structures.AuthorizeSuccessResponse, error) {
	id := os.Getenv("DISCORD_CLIENT_ID")
	secret := os.Getenv("DISCORD_CLIENT_SECRET")

	data := url.Values{}

	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirect)
	data.Set("client_id", id)
	data.Set("client_secret", secret)

	auth, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(data.Encode()))

	if err != nil {
		logger.Errorf("Failed to initialize request while getting authorization tokens from code: %v", err)
		return nil, err
	}

	auth.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := request.Do(auth)

	if err != nil {
		logger.Errorf("Failed to exchange code for authorization token: %v", err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Errorf("Failed to read body while getting authorization tokens from code: %v", err)
		return nil, err
	}

	if res.StatusCode == 400 {
		response := Structures.AuthorizeErrorResponse{}
		err := json.Unmarshal([]byte(body), &response)

		if err != nil {
			logger.Errorf("Failed to unmarshall body while getting authorization tokens from code: %v", err)
			return nil, err
		}

		return nil, errors.New(response.ErrorDescription)
	}

	response := Structures.AuthorizeSuccessResponse{}
	err = json.Unmarshal([]byte(body), &response)

	if err != nil {
		logger.Errorf("Failed to unmarshall body while getting authorization tokens from code: %v", err)
		return nil, err
	}

	return &response, nil
}

func GetDiscordUserFromAuth(auth string) (*Structures.DiscordUser, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)

	if err != nil {
		logger.Errorf("Failed to initialize request while getting authorization tokens from code: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+auth)
	res, err := request.Do(req)

	if err != nil {
		logger.Errorf("Failed to exchange code for authorization token: %v", err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Errorf("Failed to read body while getting authorization tokens from code: %v", err)
		return nil, err
	}

	if res.StatusCode == 400 {
		response := Structures.AuthorizeErrorResponse{}
		err := json.Unmarshal([]byte(body), &response)

		if err != nil {
			logger.Errorf("Failed to unmarshall body while getting authorization tokens from code: %v", err)
			return nil, err
		}

		return nil, errors.New(response.ErrorDescription)
	}

	response := Structures.DiscordUser{}
	err = json.Unmarshal([]byte(body), &response)

	if err != nil {
		logger.Errorf("Failed to unmarshall body while getting authorization tokens from code: %v", err)
		return nil, err
	}

	return &response, nil
}

func RevokeAuthorization(authorization string, redirect string) error {
	id := os.Getenv("DISCORD_CLIENT_ID")
	secret := os.Getenv("DISCORD_CLIENT_SECRET")

	data := url.Values{}

	data.Set("token_type_hint", "access_token")
	data.Set("token", authorization)
	data.Set("redirect_uri", redirect)
	data.Set("client_id", id)
	data.Set("client_secret", secret)

	auth, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token/revoke", strings.NewReader(data.Encode()))

	if err != nil {
		logger.Errorf("Failed to initialize request while revoking authorization token: %v", err)
		return err
	}

	auth.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := request.Do(auth)

	if err != nil {
		logger.Errorf("Failed to revoke authorization token: %v", err)
		return err
	}

	if res.StatusCode == 200 {
		return nil
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Errorf("Failed to read body while getting authorization tokens from code: %v", err)
		return err
	}

	logger.Infof("%v", string(body))

	message := fmt.Sprintf("Failed to revoke authorization token with status %v", res.StatusCode)

	return errors.New(message)
}
