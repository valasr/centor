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

type KiveServerInterface interface {
	Sync(*PublishRequest, string)
}

type KVMapList struct {
	Data map[string]PublishRequest
	DC   string
}
type KiveDB struct {
	DataMap       map[string]KVMapList `json:"db"`
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
	Namespace   string    `json:"namespace"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Init(serverHandler KiveServerInterface) error {
	db = &KiveDB{
		DataMap:       make(map[string]KVMapList),
		ServerHandler: serverHandler,
	}
	return nil
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

	return nil
}
func Del(ns, key string, ts int64) (*PublishRequest, error) {
	db.m.Lock()
	defer db.m.Unlock()
	var kv PublishRequest
	var n KVMapList
	var ok bool
	if n, ok = db.DataMap[ns]; ok {
		if kv, ok = n.Data[key]; ok {
			currentTs := kv.PublishDate.Unix()
			if ts < currentTs {
				return nil, fmt.Errorf("your opration is outdated : delete %s", key)
			}

			kv.Action = "delete"
			n.Data[key] = kv
		}
		delete(db.DataMap[ns].Data, key)
	} else {
		return nil, fmt.Errorf("namespace %s is not", ns)
	}

	return &kv, nil
}

func Put(dc, ns, key, value string, ts int64) (*PublishRequest, error) {
	db.m.Lock()
	defer db.m.Unlock()
	id := generateHash(key)
	var n KVMapList
	var ok bool

	if n, ok = db.DataMap[ns]; !ok {
		db.DataMap[ns] = KVMapList{
			Data: map[string]PublishRequest{},
			DC:   dc,
		}
	}

	if kv, ok := n.Data[key]; ok {
		currentTs := kv.PublishDate.Unix()
		if ts < currentTs {
			return nil, fmt.Errorf("your opration is outdated : update %s", key)
		}
	}

	kv := PublishRequest{
		Id:          id,
		PublishDate: time.Now(),
		Release:     1,
		Record: KV{
			Key:   key,
			Value: value,
		},
		Action:    "add",
		Namespace: ns,
	}
	db.DataMap[ns].Data[key] = kv

	return &kv, nil
}

func Sync(pr *PublishRequest) error {
	if db.ServerHandler == nil {
		return fmt.Errorf("server handler is nil")
	}
	db.ServerHandler.Sync(pr, "")
	return nil
}

func Get(ns, key string) (any, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	if n, ok := db.DataMap[ns]; ok {
		if kv, ok := n.Data[key]; ok {
			return kv.Record.Value, nil
		}
	} else {
		return nil, fmt.Errorf("namespace %s is not", ns)
	}

	return nil, fmt.Errorf("key not found: %s", key)
}
