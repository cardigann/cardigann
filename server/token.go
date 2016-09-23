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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"s":   t.Site,
		"l":   t.Link,
		"nbf": time.Now().Unix(),
	})

	return token.SignedString(sharedKey)
}

func decodeToken(ts string, sharedKey []byte) (*token, error) {
	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return sharedKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println(claims["foo"], claims["nbf"])
	}

	return nil, errors.New("not implemented")
}
