package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
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
