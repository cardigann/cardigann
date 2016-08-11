package server

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

func jsonError(w http.ResponseWriter, errStr string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": errStr})
}

func (h *handler) postAuthHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Passphrase string `json:"passphrase"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if h.Params.Passphrase != req.Passphrase {
		log.Printf("Client failed to authenticate")
		jsonError(w, "Invalid passphrase", http.StatusUnauthorized)
		return
	}

	log.Printf("Client successfully authenticated")
	k, err := h.sharedKey()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp = struct {
		Token string `json:"token"`
	}{
		fmt.Sprintf("%x", k),
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) sharedKey() ([]byte, error) {
	var b []byte

	switch {
	case h.Params.APIKey != nil:
		b = h.Params.APIKey
	case h.Params.Passphrase != "":
		hash := sha1.Sum([]byte(h.Params.Passphrase))
		b = hash[0:16]
	default:
		b = make([]byte, 16)
		for i := range b {
			b[i] = byte(rand.Intn(256))
		}
	}
	return b, nil
}

func (h *handler) checkAPIKey(s string) (result bool) {
	k, err := h.sharedKey()
	if err != nil {
		return false
	}
	if s == fmt.Sprintf("%x", k) {
		return true
	}
	log.Printf("Incorrect api key %q, expected %x", s, k)
	return false
}

func (h *handler) checkRequestAuthorized(r *http.Request) bool {
	if auth := r.Header.Get("Authorization"); auth != "" {
		log.Printf("Checking Authorization header")
		return h.checkAPIKey(strings.TrimPrefix(auth, "apitoken "))
	} else if apiKey := r.URL.Query().Get("apikey"); apiKey != "" {
		log.Printf("Checking apikey query string parameter")
		return h.checkAPIKey(apiKey)
	}
	return false
}
