package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/indexer"
	"github.com/gorilla/mux"
)

type indexerFeedsView struct {
	Torznab string `json:"torznab"`
}

type indexerView struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Enabled bool             `json:"enabled"`
	Feeds   indexerFeedsView `json:"feeds"`
}

type indexerViewByName []indexerView

func (slice indexerViewByName) Len() int {
	return len(slice)
}

func (slice indexerViewByName) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice indexerViewByName) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (h *handler) loadIndexerViews(baseURL string) ([]indexerView, error) {
	defs, err := indexer.ListDefinitions()
	if err != nil {
		return nil, err
	}

	reply := []indexerView{}
	for _, indexerID := range defs {
		i, err := h.lookupIndexer(indexerID)
		if err == indexer.ErrUnknownIndexer {
			log.Printf("Unknown indexer %q in configuration", indexerID)
			continue
		} else if err != nil {
			return nil, err
		}

		info := i.Info()
		reply = append(reply, indexerView{
			ID:      info.ID,
			Name:    info.Title,
			Enabled: config.IsSectionEnabled(info.ID, h.Params.Config),
			Feeds: indexerFeedsView{
				Torznab: fmt.Sprintf("%storznab/%s", baseURL, info.ID),
			},
		})
	}

	sort.Sort(indexerViewByName(reply))

	return reply, nil
}

func (h *handler) getIndexersHandler(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequestAuthorized(r) {
		jsonError(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	base, err := h.baseURL(r, "/")
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	views, err := h.loadIndexerViews(base.String())
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(views); err != nil {
		panic(err)
	}
}

func (h *handler) getIndexersConfigHandler(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequestAuthorized(r) {
		jsonError(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	indexerID := params["indexer"]

	config, err := h.Params.Config.Section(indexerID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if config == nil {
		log.
			WithFields(logrus.Fields{"indexer": indexerID}).
			Debugf("No config found for indexer")

		config = map[string]string{
			"enabled": "false",
		}
	}

	if _, ok := config["url"]; !ok {
		log.
			WithFields(logrus.Fields{"indexer": indexerID}).
			Debugf("No url found for indexer")

		i, err := h.lookupIndexer(indexerID)
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		log.WithFields(logrus.Fields{"info": i.Info()}).Debug("Loaded indexer info")
		config["url"] = i.Info().Link
	}

	log.
		WithFields(logrus.Fields{"indexer": indexerID, "config": config}).
		Debugf("Getting indexer config")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		panic(err)
	}
}

func (h *handler) patchIndexersConfigHandler(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequestAuthorized(r) {
		jsonError(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	indexerID := params["indexer"]

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	for k, v := range req {
		if err := h.Params.Config.Set(indexerID, k, v); err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (h *handler) getIndexerTestHandler(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequestAuthorized(r) {
		jsonError(w, "Not Authorized", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	indexerID := params["indexer"]

	i, err := h.lookupIndexer(indexerID)
	if err != nil {
		log.WithError(err).Error(err)
		jsonError(w, "Indexer not Found", http.StatusNotFound)
		return
	}

	if err == nil {
		tester := indexer.Tester{Runner: i}
		if err = tester.Test(); err != nil {
			log.WithError(err).Error("Test failed")
		}
	}

	var resp = struct {
		OK    bool   `json:"ok"`
		Error string `json:"error,omitempty"`
	}{}

	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.OK = true
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func (h *handler) patchIndexersHandler(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequestAuthorized(r) {
		jsonError(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	var iv indexerView
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &iv); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}
