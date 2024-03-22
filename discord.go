package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GetAuthorizationFromCode(code string, redirect string) (AuthorizeSuccessResponse, error) {
	logger.Info(code, redirect)

	id := env["DISCORD_CLIENT_ID"]
	secret := env["DISCORD_CLIENT_SECRET"]

	data := url.Values{}

	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirect)
	data.Set("client_id", id)
	data.Set("client_secret", secret)

	path := "https://discord.com/api/v10/oauth2/token"
	auth, err := http.NewRequest("POST", path, strings.NewReader(data.Encode()))

	auth.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		logger.Errorf("Failed to initialize request while getting authorization tokens from code: %v", err)
		return AuthorizeSuccessResponse{}, err
	}

	res, err := request.Do(auth)

	if err != nil {
		logger.Errorf("Failed to exchange code for authorization token: %v", err)
		return AuthorizeSuccessResponse{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Errorf("Failed to read body: %v", err)
		return AuthorizeSuccessResponse{}, err
	}

	if res.StatusCode == 400 {
		response := AuthorizeErrorResponse{}
		err := json.Unmarshal([]byte(body), &response)

		if err != nil {
			logger.Errorf("Failed to unmarshall body: %v", err)
			return AuthorizeSuccessResponse{}, err
		}

		return AuthorizeSuccessResponse{}, errors.New(response.ErrorDescription)
	}

	response := AuthorizeSuccessResponse{}
	err = json.Unmarshal([]byte(body), &response)

	if err != nil {
		logger.Errorf("Failed to unmarshall body: %v", err)
		return AuthorizeSuccessResponse{}, err
	}

	return response, nil
}
