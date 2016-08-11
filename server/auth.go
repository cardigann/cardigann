package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
