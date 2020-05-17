package store

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func (col *Collection) Put(value []byte) []byte {
	col.lastPos.mutex.Lock()
	defer col.lastPos.mutex.Unlock()

	key := uuid.New().String()
	fmt.Println("Put : key :", key)

	wNoOfBytes, err := col.writer.Write(value)
	if err != nil {
		panic(err)
	}
	col.kv[key] = valuePostion{
		Offset: col.lastPos.Pos.Offset,
		Length: wNoOfBytes,
	}
	col.lastPos.Pos.Offset = col.lastPos.Pos.Offset + int64(wNoOfBytes)
	col.lastPos.Pos.Length = wNoOfBytes

	col.SyncIndexToFile()

	return []byte(fmt.Sprintf(`{"id":"%s"}`, key))
}

func (col *Collection) Get(key string) []byte {
	col.lastPos.mutex.Lock()
	rNoOfBytes := col.kv[key].Length
	offset := col.kv[key].Offset
	col.lastPos.mutex.Unlock()

	reacolytes := make([]byte, rNoOfBytes, rNoOfBytes)
	_, err := col.reader.ReadAt(reacolytes, offset)
	if err != nil {
		panic(err)
	}
	return reacolytes
}

func (col *Collection) Delete(key string) []byte {
	col.lastPos.mutex.Lock()
	delete(col.kv, key)
	col.lastPos.mutex.Unlock()
	col.SyncIndexToFile()
	return []byte(fmt.Sprintf(`{%s}`, "ok"))
}

func (col *Collection) GetAll() []byte {
	col.lastPos.mutex.Lock()
	defer col.lastPos.mutex.Unlock()

	var values string
	for k, _ := range col.kv {
		rNoOfBytes := col.kv[k].Length
		offset := col.kv[k].Offset

		readColBytes := make([]byte, rNoOfBytes, rNoOfBytes)
		_, err := col.reader.ReadAt(readColBytes, offset)
		if err != nil {
			fmt.Println(err)
		}
		values = values + fmt.Sprintf(`{"%s":"%s"},`, k, string(readColBytes))
	}
	if len(values) > 0 {
		values = strings.TrimRight(values, ",")
	}
	collection := fmt.Sprintf(`{"%s": [%s]}`, col.name, values)
	return []byte(collection)
}
