package main

import (
	//"encoding/json"
	//"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wilke/RESTframe/CollectionJSON"
	"github.com/wilke/RESTframe/ShockClient"
	"log"
	"net/http"
	"net/url"
	//"strconv"
)

type Item CollectionJSON.Item
type Collection CollectionJSON.Collection

var myURL url.URL
var baseURL string

func init() {
	myURL.Host = "http://localhost:8000"
	baseURL = myURL.Host

	// i = new(Item)
	// 	c = new(Frame.Collection)
	// 	fmt.Printf("%+v\n", c)
	// 	fmt.Printf("%s\n", "Test")
}

func main() {

	fmt.Printf("%s\n", "Starting Server")

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.

	r.HandleFunc("/", BaseHandler)
	r.HandleFunc("/collection", CollectionHandler)
	r.HandleFunc("/shock", ShockHandler)
	r.HandleFunc("/shock/test", TestShockHandler)
	r.HandleFunc("/shock/test/{url}", TestShockHandler)
	r.HandleFunc("/shock/test/{url:.+}", TestShockHandler)
	// 	r.HandleFunc("/experiment/{id:[a-zA-Z]*}", ExperimentHandler).Name("experiment")
	// 	r.HandleFunc("/search", SearchHandler)
	// 	r.HandleFunc("/search/{path:.+}", SearchHandler)
	// 	r.HandleFunc("/register", RegisterHandler)
	// 	r.HandleFunc("/register/{path:[a-z+]+}", RegisterHandler)
	// 	r.HandleFunc("/upload", UploadHandler)
	// 	r.HandleFunc("/download", DownloadHandler)
	// 	r.HandleFunc("/transfer", TransferHandler)
	// 	r.HandleFunc("/transfer/{id}", SearchHandler)
	// 	r.HandleFunc("/test", GetExperimentHandler)

	// Bind to a port and pass our router in
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Printf("%s\n", "Started Server at port 8000")
}

func BaseHandler(w http.ResponseWriter, r *http.Request) {

	c := new(CollectionJSON.CollectionJSON)

	jb, err := c.ToJson()

	// Send json
	if err != nil {
		println(jb)
		w.Write([]byte(err.Error()))
		http.Error(w, err.Error(), 500)
	} else {
		fmt.Printf("%s\n", jb)
		fmt.Printf("%+v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jb))

	}

}

// CollectionJSON struct
func CollectionHandler(w http.ResponseWriter, r *http.Request) {

	c := new(CollectionJSON.CollectionJSON)

	jb, err := c.ToJson()

	// Send json
	if err != nil {
		println(jb)
		w.Write([]byte(err.Error()))
		http.Error(w, err.Error(), 500)
	} else {
		fmt.Printf("%s\n", jb)
		fmt.Printf("%+v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jb))

	}

}

// Shock struct
func ShockHandler(w http.ResponseWriter, r *http.Request) {

	c := new(ShockClient.Collection)

	jb, err := c.ToJson()

	// Send json
	if err != nil {
		println(jb)
		w.Write([]byte(err.Error()))
		http.Error(w, err.Error(), 500)
	} else {
		fmt.Printf("%s\n", jb)
		fmt.Printf("%+v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jb))

	}

}

func TestShockHandler(w http.ResponseWriter, r *http.Request) {

	// Retrieve path and query parameters for experiment resource
	vars := mux.Vars(r)
	url := vars["url"]
	status := 500

	// initialize r.Form for query parameters
	r.ParseForm()

	var err error
	var collection ShockClient.Collection
	client := new(ShockClient.Client)

	if url != "" {
		client.URL = url

		fmt.Println("GET url")
		collection, status, err = client.Get("http://" + url)
		collection.Status = status

		if err != nil {
			//client.SendError(w, getERR, 500
			collection.Error = err.Error()
		}

	}

	jb, jsonERR := collection.ToJson()

	// Send json
	if jsonERR != nil {
		println(jb)
		w.Write([]byte(jsonERR.Error()))
		http.Error(w, jsonERR.Error(), 500)
	} else {
		//fmt.Printf("%s\n", jb)
		//fmt.Printf("%+v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jb))

	}

}
