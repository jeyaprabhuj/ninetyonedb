package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"ninetyonedb/store"
	"ninetyonedb/store/http/helpers"
	"ninetyonedb/store/http/response"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var KVServer *store.KVServer
var PORT string

func init() {
	KVServer = store.CreateKVServer()
	for _, storage := range KVServer.Storages {
		storageRoutePath := "/storages/" + storage.Name + "/"
		http.HandleFunc(storageRoutePath, storageHandler)
		http.HandleFunc(storageRoutePath+"collections/", IntializeCollectionHandler)
		for _, collection := range storage.GetCollections() {
			collectionRoutePath := "collections/" + collection + "/"
			http.HandleFunc(storageRoutePath+collectionRoutePath, collectionhandler)
		}
	}
}

var collectionhandler = func(w http.ResponseWriter, r *http.Request) {
	storage, collection, documentId := helpers.DestructurePath(r.URL.Path)
	collectionPtr := KVServer.GetStorage(storage).GetCollection(collection)

	if r.Method == http.MethodGet {
		if documentId != "" {
			helpers.CreateJSONResponse(w, collectionPtr.Get(documentId))

		} else {
			helpers.CreateJSONResponse(w, collectionPtr.GetAll())
		}
	} else if r.Method == http.MethodPut {
		buffer, _ := ioutil.ReadAll(r.Body)
		var raw map[string]interface{}
		if err := json.Unmarshal(buffer, &raw); err != nil {
			helpers.CreateJSONResponse(w, helpers.CreateErrorResponse(404, err.Error()))
			fmt.Println(err)
			return
		}
		helpers.CreateJSONResponse(w, collectionPtr.Put(buffer))

	} else if r.Method == http.MethodDelete {
		buffer, _ := ioutil.ReadAll(r.Body)
		var raw map[string]interface{}
		if err := json.Unmarshal(buffer, &raw); err != nil {
			helpers.CreateJSONResponse(w, helpers.CreateErrorResponse(404, err.Error()))
			fmt.Println(err)
		}
		var id string
		for k, v := range raw {
			if "id" == strings.ToLower(fmt.Sprintf("%s", k)) {
				id = fmt.Sprintf("%s", v)
			}
		}
		if id != "" {
			helpers.CreateJSONResponse(w, collectionPtr.Delete(id))
		}
	}
}

var storageHandler = func(w http.ResponseWriter, r *http.Request) {
	// storage, collection, documentId := helpers.DestructurePath(r.URL.Path)
	// storagePtr := KVServer.GetStorage(storage)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("Database Handler"))
	if r.Method == http.MethodGet {

	}
}

var IntializeCollectionHandler = func(w http.ResponseWriter, r *http.Request) {
	storage, _, _ := helpers.DestructurePath(r.URL.Path)
	storagePtr := KVServer.GetStorage(storage)
	storageRoutePath := "/storages/" + storagePtr.Name + "/"
	if r.Method == http.MethodGet {
		bCollections, err := json.Marshal(storagePtr.GetCollections())
		if err != nil {
			helpers.CreateJSONResponse(w, helpers.CreateErrorResponse(404, err.Error()))
			fmt.Println(err)
		}
		helpers.CreateJSONResponse(w, bCollections)
	} else if r.Method == http.MethodPut {
		buffer, _ := ioutil.ReadAll(r.Body)
		var raw map[string]interface{}
		if err := json.Unmarshal(buffer, &raw); err != nil {
			helpers.CreateJSONResponse(w, helpers.CreateErrorResponse(404, err.Error()))
			fmt.Println(err)
		}
		collection, _ := raw["name"].(string)
		storagePtr.CreateCollection(collection)
		collectionRoutePath := "collections/" + collection + "/"
		http.HandleFunc(storageRoutePath+collectionRoutePath, collectionhandler)
		helpers.CreateJSONResponse(w, []byte(`{"ok"}`))
	}
}

var InitializeStoragesHandlers = func(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		bStorage, err := json.Marshal(KVServer.GetStorageNames())
		if err != nil {
			helpers.CreateJSONResponse(w, helpers.CreateErrorResponse(404, err.Error()))
			fmt.Println(err)
		}
		helpers.CreateJSONResponse(w, bStorage)
	} else if r.Method == http.MethodPut {
		buffer, _ := ioutil.ReadAll(r.Body)
		var raw map[string]interface{}
		if err := json.Unmarshal(buffer, &raw); err != nil {
			helpers.CreateJSONResponse(w, helpers.CreateErrorResponse(404, err.Error()))
			fmt.Println(err)
		}
		storage, _ := raw["name"].(string)
		KVServer.CreateKeyValueStore(storage)
		http.HandleFunc("/storages/"+storage+"/", storageHandler)
		http.HandleFunc("/storages/"+storage+"/collections/", IntializeCollectionHandler)
		helpers.CreateJSONResponse(w, []byte(`{"ok"}`))
	}
}

var ServerInfoHandler = func(w http.ResponseWriter, r *http.Request) {
	pr, pw := io.Pipe()

	info := &response.ServerInfo{
		"v1.0beta",
		"Json Storage server",
		"jeyaprabhu.j@gmail.com => Varshini Ambroz Industries",
		KVServer.GetStorageNames(),
		fmt.Sprintf("http://localhost:%s/storages", PORT),
	}
	go func() {
		defer pw.Close()
		json.NewEncoder(pw).Encode(info)
	}()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.Copy(w, pr)
}

func checkPort() error {
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Println("Port error : ", err)
	}
	defer listener.Close()
	return err
}
func main() {

	var srv *http.Server
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	if !SetPort() {
		c <- syscall.SIGINT
	}

	http.HandleFunc("/", ServerInfoHandler)
	http.HandleFunc("/storages/", InitializeStoragesHandlers)
	srv = &http.Server{Addr: ":" + PORT}

	var cleanup = func() {
		signal.Stop(c)
		cancel()
		fmt.Println("Clean up in progress .....Please wait ")
		srv.Shutdown(ctx)
		KVServer.Close()
	}
	defer func() {
		fmt.Println("Deferred call")
		cleanup()
	}()

	go func() {
		select {
		case <-c:
			cancel()
			fmt.Println("Shutting Down....... ")
			cleanup()
		case <-ctx.Done():
			fmt.Println("Context done.....")
		}
	}()

	fmt.Printf("Servig @: http://localhost:%s/storages \n", PORT)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}

func SetPort() bool {
	PORT = os.Getenv("KVPORT")
	fmt.Println("PORT from os environment : ", PORT)
	if PORT != "" {
		if checkPort() != nil {
			return false
		}
	} else {
		PORT = "8001"
	}

	return true
}
