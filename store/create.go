package store

import (
	"os"
	"time"
)

func CreateKVServer() *KVServer {
	kvServer := &KVServer{}
	kvServer.Storages = make(map[string]*Storage)
	kvServer.storageFile = "./storages.kvserver"
	kvServer.ReadStoragesFromFile()
	go func() {
		for {
			time.Sleep(60 * time.Second)
			kvServer.SyncStoragesToFile()
		}
	}()
	return kvServer
}

func (kvServer *KVServer) CreateKeyValueStore(Name string) *Storage {
	if _, err := os.Stat(Name); os.IsNotExist(err) {
		os.Mkdir(Name, 0777)
	}

	db := &Storage{}
	db.Name = Name
	db.collections = make(map[string]*Collection)
	db.colfile = "./" + db.Name + "/collections"
	db.ReadCollectionsFromFile()
	go func() {
		for {
			time.Sleep(60 * time.Second)
			db.SyncCollectionsToFile()
		}
	}()
	return db
}

func (s *Storage) CreateCollection(name string) *Collection {
	col := &Collection{}
	col.datafolder = "./" + s.Name + "/" + name + "/"
	if _, err := os.Stat(col.datafolder); os.IsNotExist(err) {
		os.Mkdir(col.datafolder, 0777)
	}

	dataFile := col.datafolder + "data"
	col.name = name
	col.writer, _ = os.OpenFile(dataFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	col.reader, _ = os.OpenFile(dataFile, os.O_RDONLY, 0777)
	col.kv = make(map[string]valuePostion)
	col.InitializeCollection()
	s.collections[name] = col
	return col
}
