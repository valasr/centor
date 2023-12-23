package kive

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	dbPath = "/var/lib/centor/kive/db"
)

var db *KiveDB

type KiveDB struct {
	data map[string]PublishRequest
	m    sync.RWMutex
}

type PublishRequest struct {
	Id          string    `json:"id"`
	RequestDate time.Time `json:"request_date"`
	PublishDate time.Time `json:"publish_date"`
	Release     int       `json:"release"` // 1 : in-hard, 2 : in-cluster, 4 :
	Record      KV        `json:"record"`
}

type KV struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func LoadDB() error {
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		log.Println("db is empty")
		return nil
	}

	err = json.Unmarshal(data, db)
	if err != nil {
		return err
	}

	fmt.Printf("data : %+v", db.data)

	return nil
}

func (k *KiveDB) Set(key string, value any) {
	k.m.Lock()
	defer k.m.Unlock()
	id := generateHash(key)

	db.data[id] = PublishRequest{
		Id:          id,
		PublishDate: time.Now(),
		Release:     1,
		Record: KV{
			Key:   key,
			Value: value,
		},
	}

	jsonData, err := json.Marshal(db.data)
	if err != nil {
		log.Fatal("error in marshalling : ", err)
	}
	err = os.WriteFile(dbPath, jsonData, 0655)
	if err != nil {
		log.Fatal("error in writing : ", err)
	}

	if r, ok := db.data[id]; ok {
		r.Release += 2
		db.data[id] = r
	}
}

func (k *KiveDB) Get(key string) (any, error) {
	k.m.RLock()
	defer k.m.RUnlock()

	if r, ok := db.data[key]; ok {
		return r.Record.Value, nil
	}
	return nil, fmt.Errorf("key not found: %s", key)
}
