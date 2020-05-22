package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpHandlers struct {
	indexer *indexer
}

func newHTTPHandlers(i *indexer) *httpHandlers {
	return &httpHandlers{
		indexer: i,
	}
}

func (h *httpHandlers) root(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.search(w, r)
		return
	case "POST":
		h.addItem(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h *httpHandlers) addItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	err = h.indexer.analyzeAndIndex(string(bytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}
}

func (h *httpHandlers) search(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("q")

	docs, err := h.indexer.search(term)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	err = json.NewEncoder(w).Encode(docs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

}
