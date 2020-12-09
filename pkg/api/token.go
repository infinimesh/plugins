package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type TokenHandler interface {
	RefreshToken() error
	GetToken() string
}

type tokenHandlerImpl struct {
	client    *http.Client
	username  string
	password  string
	apiUrl    string
	lastToken string
}

// NewTokenHandler returns a new TokenHandler
//
// Tokens are synchronously retrieved on init, and refreshed automatically per
// tick duration
//
// If a token retrieval attempt fails, GetToken returns the previous token
func NewTokenHandler(username, password, apiUrl string, refreshInterval time.Duration) TokenHandler {
	h := &tokenHandlerImpl{
		client:   &http.Client{},
		username: username,
		password: password,
		apiUrl:   apiUrl,
	}
	if err := h.RefreshToken(); err != nil {
		log.Printf("error retrieving token: %s\n", err)
	}
	go func() {
		for range time.Tick(refreshInterval) {
			if err := h.RefreshToken(); err != nil {
				log.Printf("error retrieving token: %s\n", err)
			}
		}
	}()
	return h
}

type tokenParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenRes struct {
	Token string `json:"token"`
}

func (h *tokenHandlerImpl) RefreshToken() error {
	log.Println("refreshing account token...")
	b, _ := json.Marshal(tokenParams{
		Username: h.username,
		Password: h.password,
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, h.apiUrl+endpointAccountToken, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")

	res, err := h.client.Do(req)
	if err != nil {
		return NewFailedHTTPRequestError(http.MethodPost, endpointAccountToken, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return NewUnexpectedHTTPStatusCodeError(http.MethodPost, endpointAccountToken, res.StatusCode)
	}

	b, _ = ioutil.ReadAll(res.Body)
	tmp := tokenRes{}
	_ = json.Unmarshal(b, &tmp)
	h.lastToken = tmp.Token
	return nil
}

func (h *tokenHandlerImpl) GetToken() string {
	return h.lastToken
}
