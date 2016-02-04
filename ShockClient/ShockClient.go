package ShockClient

import (
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/wilke/webserver/CollectionJson"
	//"log"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	//"time"
)

type Client struct {
	URL  string
	User User
}

type ShockClient Client

type User struct {
	login   string
	pwd     string
	authUrl string
	token   string
}

type File struct {
	Name          string            `json:"name"`
	Size          int               `json:"size"`
	Checksum      map[string]string `json:"checksum"`
	Format        string            `json:"format"`
	Virtual       bool              `json:"virtual"`
	Virtual_parts []interface{}     `json:"virtual_parts"`
	Created_on    string            `json:"created_on"`
}

type Attributes map[string]interface{}

type Node struct {
	Id            string      `json:"ID"`
	Version       string      `json:"version"`
	File          File        `json:"file"`
	Attributes    interface{} `json:"attributes"`
	Indexes       interface{} `json:"indexes"`
	Tags          []string    `json:"tags"`
	Linkage       []string    `json:"linkage"`
	Created_on    string      `json:"created_on"`
	Last_modified string      `json:"last_modified"`
	Expiration    string      `json:"expiration"`
	Type          string      `json:"type"`
	Parts         []string    `json:"parts"`
}

type Collection struct {
	Error       string      `json:"error"`
	Data        interface{} `json:"data"`
	Status      int         `json:"status"`
	Limit       int         `json:"limit"`
	Offset      int         `json:"offset"`
	Total_count int         `json:"total_count"`
}

type NodeCollection struct {
	Collection
	Data Node `json:"data"`
}

type NodesCollection struct {
	Collection
	Data []Node `json:"data"`
}

type Frame struct{}

// Methods
func ToJson(d interface{}) ([]byte, error) {
	jb, err := json.Marshal(d)
	if err != nil {
		fmt.Printf("%+v\n", err)
	} else {
		fmt.Printf("%s\n", jb)
		//fmt.Printf("%+v\n", err)
	}
	return jb, err
}

func (c Collection) ToJson() ([]byte, error) {
	jb, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("%+v\n", err)
	} //else {
	//fmt.Printf("%s\n", jb)
	//fmt.Printf("%+v\n", err)
	//}
	return jb, err
}

// Client methods
func (c Client) GetToken() (string, error)            { return "", nil }
func (c Client) SetAuthHeader() (string, error)       { return "", nil }
func (c Client) CheckAuthHeader() (string, error)     { return "", nil }
func (c Client) Post(url string, d interface{}) error { return nil }
func (c Client) Put(url string, d interface{}) error  { return nil }

func (c Client) Get(uri string) (Collection, int, error) {

	var collection Collection
	var tmp_single NodeCollection
	var tmp_multiple NodesCollection
	var nodes []Node
	var status int

	url, err := url.Parse(uri)

	fmt.Printf("Retrieving data from %+v\n", url.String())
	resp, err := http.Get(url.String())

	if err != nil {
		fmt.Printf("%+v\n", err)
		status = resp.StatusCode
		return collection, status, err
	} else {
		status = resp.StatusCode
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &tmp_single)
	if err != nil {
		err2 := json.Unmarshal(body, &tmp_multiple)
		if err2 != nil {
			fmt.Printf("Something very wrong. Can't create collection from json:\n1: %v \n 2: %v", err.Error, err2.Error)
			status = 500
			return collection, status, errors.New(strings.Join([]string{err.Error(), err2.Error()}, " "))
		}
		collection.Error = tmp_multiple.Error
		collection.Status = tmp_multiple.Status
		collection.Limit = tmp_multiple.Limit
		collection.Offset = tmp_multiple.Offset
		collection.Total_count = tmp_multiple.Total_count
		collection.Data = tmp_multiple.Data

		fmt.Printf("Mutiple %+v\n", collection)
		return collection, status, nil
	} else {
		nodes = append(nodes, tmp_single.Data)
		collection.Error = tmp_multiple.Error
		collection.Data = nodes
		return collection, status, nil
	}

	//nodes = append(nodes, tmp.Data.([]Node)[0])
	fmt.Printf("Why %+v\n", collection)
	return collection, status, nil
}

func (c Client) SendError(w http.ResponseWriter, err error, errCode int) {

	if errCode != 0 {
		errCode = 500
	}
	http.Error(w, err.Error(), errCode)
}

func (c Client) Send(w http.ResponseWriter, d interface{}) error {

	jb, err := ToJson(d)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jb))

	return nil
}

// Helper function

func FillStruct(data map[string]interface{}, result interface{}) {
	t := reflect.ValueOf(result).Elem()
	for k, v := range data {
		val := t.FieldByName(k)
		fmt.Println("Test", k, v)
		if reflect.TypeOf(v) != nil {
			fmt.Printf("Type %s\n", reflect.TypeOf(v).Name())
			fmt.Printf("Value %+d\n", reflect.ValueOf(v))
			fmt.Printf("Field %d\n", val)
			if reflect.ValueOf(v).IsValid() {
				if val.IsValid() {
					val.Set(reflect.ValueOf(v))
				}
			}
		}
	}
}
