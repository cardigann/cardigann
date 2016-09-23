package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type token struct {
	Site string `json:"s,omitempty"`
	Link string `json:"l,omitempty"`
}

func (t *token) Encode(sharedKey []byte) (string, error) {
	j := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"s":   t.Site,
		"l":   t.Link,
		"nbf": time.Now().Unix(),
	})

	return j.SignedString(sharedKey)
}

func decodeToken(ts string, sharedKey []byte) (*token, error) {
	j, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return sharedKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := j.Claims.(jwt.MapClaims)
	if !ok || !j.Valid {
		return nil, errors.New("Invalid token")
	}

	return &token{claims["s"].(string), claims["l"].(string)}, nil
}
