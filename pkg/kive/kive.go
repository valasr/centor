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

func init() {
	db = &KiveDB{
		Data: make(map[string]PublishRequest),
	}
}

type KiveServerInterface interface {
	Sync(PublishRequest)
}
type KiveDB struct {
	Data          map[string]PublishRequest `json:"data"`
	ServerHandler KiveServerInterface
	m             sync.RWMutex
}

type PublishRequest struct {
	Id          string    `json:"id"`
	RequestDate time.Time `json:"request_date"`
	PublishDate time.Time `json:"publish_date"`
	Release     int       `json:"release"` // 1 : in-hard, 2 : in-cluster, 4 :
	Record      KV        `json:"record"`
	Action      string    `json:"action"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

	fmt.Printf("data : %+v", db.Data)

	return nil
}
func Del(id string) (*PublishRequest, error) {
	db.m.Lock()
	defer db.m.Unlock()
	var kv PublishRequest
	var ok bool
	if kv, ok = db.Data[id]; ok {
		kv.Action = "delete"
	}
	delete(db.Data, id)
	return &kv, nil
}

func Put(key, value string) (*PublishRequest, error) {
	db.m.Lock()
	defer db.m.Unlock()
	id := generateHash(key)

	kv := PublishRequest{
		Id:          id,
		PublishDate: time.Now(),
		Release:     1,
		Record: KV{
			Key:   key,
			Value: value,
		},
		Action: "add",
	}
	db.Data[key] = kv

	// jsonData, err := json.Marshal(db.Data)
	// if err != nil {
	// 	return fmt.Errorf("error in marshalling : %s", err)
	// }
	// err = os.WriteFile(dbPath, jsonData, 0655)
	// if err != nil {
	// 	return fmt.Errorf("error in writing : %s", err)
	// }

	// if r, ok := db.Data[id]; ok {
	// 	r.Release += 2
	// 	db.Data[id] = r
	// }

	return &kv, nil

}

func Get(key string) (any, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	if r, ok := db.Data[key]; ok {
		return r.Record.Value, nil
	}

	return nil, fmt.Errorf("key not found: %s", key)
}
