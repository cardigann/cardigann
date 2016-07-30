package server

import (
	"encoding/json"
	"time"

	"github.com/dvsekhvalnov/jose2go"
)

type token struct {
	Site string `json:"s,omitempty"`
	Link string `json:"l,omitempty"`
}

func (t *token) Encode(sharedKey []byte) (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return jose.Encrypt(string(b), jose.A128GCMKW, jose.A128GCM, sharedKey,
		jose.Headers(map[string]interface{}{
			"c": time.Now().String(),
		}))
}

func decodeToken(ts string, sharedKey []byte) (*token, error) {
	payload, _, err := jose.Decode(ts, sharedKey)
	if err != nil {
		return nil, err
	}

	var t token
	return &t, json.Unmarshal([]byte(payload), &t)
}
