package store

import (
	"os"
	"sync"
)

type KVServer struct {
	Storages    map[string]*Storage
	storageFile string
	mutex       sync.Mutex
}

type Storage struct {
	Name    string
	collections map[string]*Collection
	colfile     string
	mutex       sync.Mutex
}

type Collection struct {
	name       string
	writer     *os.File
	reader     *os.File
	lastPos    lastValuePosition
	kv         map[string]valuePostion
	datafolder string
	mutex      sync.Mutex
	syncindex  int
}

type valuePostion struct {
	Offset int64
	Length int
}

type lastValuePosition struct {
	Pos   valuePostion
	mutex sync.Mutex
}
