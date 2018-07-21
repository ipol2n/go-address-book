package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database
var collection *mgo.Collection

// Address ...
type Address struct {
	ID      bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string        `json:"name"`
	Address string        `json:"address"`
	Tel     string        `json:"tel"`
}

// Connect ...
func Connect() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB("address-book")
	collection = db.C("address")
}

func init() {
	Connect()
}

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello World")
	// })
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var address Address
		err = json.Unmarshal(b, &address)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		collection.Insert(&address)
		output, err := json.Marshal(address)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)
	})
	http.HandleFunc("/record", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			var addresses []Address
			err := collection.Find(bson.M{}).All(&addresses)
			if err == nil {
				fmt.Println(addresses)
				s, _ := json.Marshal(addresses)
				fmt.Fprintf(w, string(s))
			} else {
				fmt.Println(err)
			}
		}
	})
	http.HandleFunc("/record/", func(w http.ResponseWriter, r *http.Request) {
		ID := strings.TrimPrefix(r.URL.Path, "/record/")
		fmt.Println("ID", ID)
		if r.Method == "GET" {
			var address Address
			err := collection.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).One(&address)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			s, _ := json.Marshal(address)
			fmt.Fprintf(w, string(s))
		} else if r.Method == "PUT" {

			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var address Address
			err = json.Unmarshal(b, &address)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			err = collection.Update(bson.M{"_id": bson.ObjectIdHex(ID)}, bson.M{"$set": bson.M{"name": address.Name, "address": address.Address, "tel": address.Tel}})
		} else if r.Method == "DELETE" {
			collection.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})
		}
	})
	http.ListenAndServe(":3000", nil)
	fmt.Println("Hello")
}
