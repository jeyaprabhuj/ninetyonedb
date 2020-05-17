package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func (kvServer *KVServer) GetStorageNames() []string {
	storages := make([]string, 0)

	for k, _ := range kvServer.Storages {
		storages = append(storages, k)
	}
	return storages
}

func (kvServer *KVServer) GetStorages() []*Storage {
	storages := make([]*Storage, 0)

	for _, v := range kvServer.Storages {
		storages = append(storages, v)
	}
	return storages
}

func (kvServer *KVServer) GetStorage(name string) *Storage {
	return kvServer.Storages[name]
}

func (kvServer *KVServer) SyncStoragesToFile() {
	kvServer.mutex.Lock()
	defer kvServer.mutex.Unlock()
	os.Remove(kvServer.storageFile)
	var storages string
	for k, _ := range kvServer.Storages {
		storages = k + ","
	}
	ioutil.WriteFile(kvServer.storageFile, []byte(storages), 0777)
	fmt.Println("collection file sync ....")
}

func (kvServer *KVServer) ReadStoragesFromFile() {
	kvServer.mutex.Lock()
	defer kvServer.mutex.Unlock()

	bStorages, err := ioutil.ReadFile(kvServer.storageFile)
	if err != nil {
		fmt.Println(err)
	} else {
		strcols := strings.Split(string(bStorages), ",")
		for _, storage := range strcols {
			if len(storage) > 0 {
				if kvServer.Storages[storage] != nil {
					kvServer.Storages[storage].Close()
				}
				delete(kvServer.Storages, storage)
				kvServer.Storages[storage] = kvServer.CreateKeyValueStore(storage)
			}
		}
	}
}

func (kvServer *KVServer) Close() {
	for _, storage := range kvServer.Storages {
		storage.Close()
	}
}
