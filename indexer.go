package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"sync"

	"github.com/boltdb/bolt"
)

type indexer struct {
	sync.Mutex
	nextId   uint32
	analyzer *analyzer
	db       *bolt.DB
}

func newIndexer(db *bolt.DB) *indexer {
	return &indexer{
		analyzer: &analyzer{},
		db:       db,
	}
}

func (i *indexer) getID() uint32 {
	i.Lock()
	defer i.Unlock()

	current := i.nextId
	i.nextId++
	return current
}

func (i *indexer) analyzeAndIndex(doc string) error {
	docId := i.getID()
	words := i.analyzer.splitAndLowercase(doc)

	err := i.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Data"))
		if err := b.Put([]byte(fmt.Sprintf("%d", docId)), []byte(doc)); err != nil {
			return err
		}

		b = tx.Bucket([]byte("InvertedIndex"))
		for _, word := range words {
			key := []byte(word)
			data := b.Get(key)
			updated := bytes.NewBuffer(data)
			binary.Write(updated, binary.LittleEndian, &docId)
			err := b.Put(key, updated.Bytes())
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

type doc struct {
	ID      uint32 `json:"id"`
	Content string `json:"content"`
}

func (i *indexer) search(term string) ([]doc, error) {
	term = strings.ToLower(term)
	var docIdBytes []byte
	if err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("InvertedIndex"))
		data := b.Get([]byte(term))
		docIdBytes = data
		return nil

	}); err != nil {
		return nil, err
	}

	if len(docIdBytes) == 0 {
		return nil, nil
	}

	docIds := make([]uint32, len(docIdBytes)/4)
	err := binary.Read(bytes.NewReader(docIdBytes), binary.LittleEndian, &docIds)
	if err != nil {
		return nil, err
	}

	docs := make([]doc, len(docIds))
	for i, id := range docIds {
		var content string
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Data"))
			content = string(b.Get([]byte(fmt.Sprintf("%d", id))))
			return nil
		})
		docs[i] = doc{
			ID:      id,
			Content: content,
		}

	}
	return docs, nil
}
