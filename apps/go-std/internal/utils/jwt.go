package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

func DecodeJwt(jwt string) (map[string]interface{}, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid jwt")
	}
	jsonPayload, err := base64.RawURLEncoding.WithPadding(base64.NoPadding).DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var payload map[string]interface{}
	err = json.Unmarshal(jsonPayload, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
