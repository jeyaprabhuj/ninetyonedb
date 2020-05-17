package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

func (col *Collection) SyncIndexToFile() {
	func() {
		idxFileName := col.GetDataFilePath("idx")
		col.mutex.Lock()
		defer col.mutex.Unlock()
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err := enc.Encode(col.kv)
		if err != nil {
			fmt.Println(err)
		}
		os.Remove(idxFileName)
		ioutil.WriteFile(idxFileName, buffer.Bytes(), 0777)
		fmt.Println("index file sync ....")
	}()
	col.LastPositionToFile()
}

func (col *Collection) LastPositionToFile() {
	lasposfile := col.GetDataFilePath("lasPos")
	col.mutex.Lock()
	defer col.mutex.Unlock()
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(col.lastPos.Pos)
	if err != nil {
		fmt.Println(err)
	}
	os.Remove(lasposfile)
	ioutil.WriteFile(lasposfile, buffer.Bytes(), 0777)
}
func (col *Collection) LastPositionFromFile() {
	lasposfile := col.GetDataFilePath("lasPos")
	col.mutex.Lock()
	defer col.mutex.Unlock()

	fileBytes, err := ioutil.ReadFile(lasposfile)
	if err != nil {
		fmt.Println(err)
	}
	buffer := bytes.NewBuffer(fileBytes)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(&col.lastPos.Pos)
	if err != nil {
		fmt.Println(err)
	}
}
func (col *Collection) ReadIndexFromFile() int {
	var fileBytes []byte
	var err error
	func() {
		col.mutex.Lock()
		defer col.mutex.Unlock()

		fileBytes, err = ioutil.ReadFile(col.GetDataFilePath("idx"))
		if err != nil {
			fmt.Println(err)
		}
		buffer := bytes.NewBuffer(fileBytes)
		dec := gob.NewDecoder(buffer)
		err = dec.Decode(&col.kv)
		if err != nil {
			fmt.Println(err)
		}
	}()
	col.LastPositionFromFile()
	return len(fileBytes)
}

func (col *Collection) InitializeCollection() {
	if col.ReadIndexFromFile() == 0 {
		col.SyncIndexToFile()
	}
}

func (col *Collection) Close() {

	col.SyncIndexToFile()
	col.writer.Close()
	col.reader.Close()
}

func (col *Collection) GetDataFilePath(name string) string {
	return col.datafolder + name
}
