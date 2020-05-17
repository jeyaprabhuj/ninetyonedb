package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func (s *Storage) GetCollection(name string) *Collection {
	return s.collections[name]
}

func (s *Storage) GetCollections() []string {
	cols := make([]string, 0)

	for k, _ := range s.collections {
		cols = append(cols, k)
	}
	return cols
}

func (s *Storage) Close() {
	for _, v := range s.collections {
		v.Close()
	}
}

func (s *Storage) SyncCollectionsToFile() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var cols string
	os.Remove(s.colfile)
	for k, _ := range s.collections {
		cols = k + ","
	}
	ioutil.WriteFile(s.colfile, []byte(cols), 0777)
	fmt.Println("collection file sync ....")
}

func (s *Storage) ReadCollectionsFromFile() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bcols, err := ioutil.ReadFile(s.colfile)
	if err != nil {
		fmt.Println(err)
	} else {
		strcols := strings.Split(string(bcols), ",")
		for _, col := range strcols {
			if len(col) > 0 {
				if s.collections[col] != nil {
					s.collections[col].Close()
				}
				delete(s.collections, col)
				s.collections[col] = s.CreateCollection(col)
			}
		}
	}
}
