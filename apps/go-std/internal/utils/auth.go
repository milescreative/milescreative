package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
)

var secret = []byte("NK9LELM2XGM89R5XYJ6S====")

func GenerateCSRFToken(sessionId string) (string, []byte) {

	csrfToken, err := GenerateRandomString()
	if err != nil {
		log.Fatal(err)
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(csrfToken + "." + sessionId))
	csrfTokenHMAC := mac.Sum(nil)
	return csrfToken, csrfTokenHMAC
}

func VerifyCSRFToken(token string, sessionId string, storedHMAC []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(token + "." + sessionId))
	expectedHMAC := mac.Sum(nil)
	log.Println("expectedHMAC: ", expectedHMAC)
	log.Println("storedHMAC: ", storedHMAC)
	return hmac.Equal(storedHMAC, expectedHMAC)
}

func GenerateSessionToken() (string, error) {
	//TODO : generate hash for storage and function to verify similar to csrf
	state, err := GenerateRandomStringNoPadding()
	if err != nil {
		log.Fatal(err)
	}
	return state, nil
}

func GenerateState() (string, error) {
	state, err := GenerateRandomStringNoPadding()
	if err != nil {
		log.Fatal(err)
	}
	return state, nil
}

func GenerateCodeVerifier() (string, error) {
	state, err := GenerateRandomStringNoPadding()
	if err != nil {
		log.Fatal(err)
	}
	return state, nil
}

func CreateS256CodeChallenge(codeVerifier string) string {
	codeChallenge := sha256.Sum256([]byte(codeVerifier))
	return EncodeBase64UrlNoPadding(codeChallenge[:])
}

func EncodeBasicCredentials(clientID string, clientSecret string) string {
	bytes := []byte(clientID + ":" + clientSecret)
	return base64.StdEncoding.EncodeToString(bytes)
}

type OAuth2Tokens struct {
	Data map[string]interface{}
}

func NewOAuth2Tokens(body []byte) (*OAuth2Tokens, error) {
	var t OAuth2Tokens
	if err := json.Unmarshal(body, &t.Data); err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *OAuth2Tokens) TokenType() (string, error) {
	v, ok := t.Data["token_type"].(string)
	if !ok {
		return "", errors.New("missing/invalid 'token_type'")
	}
	return v, nil
}

func (t *OAuth2Tokens) AccessToken() (string, error) {
	v, ok := t.Data["access_token"].(string)
	if !ok {
		return "", errors.New("missing/invalid 'access_token'")
	}
	return v, nil
}

func (t *OAuth2Tokens) AccessTokenExpiresInSeconds() (float64, error) {
	v, ok := t.Data["expires_in"].(float64)
	if !ok {
		return 0, errors.New("missing/invalid 'expires_in'")
	}
	return v, nil
}

func (t *OAuth2Tokens) AccessTokenExpiresAt() (time.Time, error) {
	seconds, err := t.AccessTokenExpiresInSeconds()
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().Add(time.Duration(seconds) * time.Second), nil
}

func (t *OAuth2Tokens) HasRefreshToken() bool {
	_, ok := t.Data["refresh_token"].(string)
	return ok
}

func (t *OAuth2Tokens) RefreshToken() (string, error) {
	v, ok := t.Data["refresh_token"].(string)
	if !ok {
		return "", errors.New("missing/invalid 'refresh_token'")
	}
	return v, nil
}

func (t *OAuth2Tokens) HasScopes() bool {
	_, ok := t.Data["scope"].(string)
	return ok
}

func (t *OAuth2Tokens) Scopes() ([]string, error) {
	v, ok := t.Data["scope"].(string)
	if !ok {
		return nil, errors.New("missing/invalid 'scope'")
	}
	return strings.Split(v, " "), nil
}

func (t *OAuth2Tokens) IDToken() (string, error) {
	v, ok := t.Data["id_token"].(string)
	if !ok {
		return "", errors.New("missing/invalid 'id_token'")
	}
	return v, nil
}
