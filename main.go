package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func main() {
	initBolt()

	i := newIndexer(db)
	h := newHTTPHandlers(i)
	http.HandleFunc("/", h.root)
	http.ListenAndServe(":9090", nil)
}

func initBolt() {
	boltdb, err := bolt.Open("./data/bolt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db = boltdb

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("InvertedIndex"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte("Data"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
