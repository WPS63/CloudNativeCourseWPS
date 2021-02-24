package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{"shoes": 50, "socks": 5} //create an instance of type database
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/delete", db.delete)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

var mutex = &sync.RWMutex{}

type dollars float64 //declare type dollars

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) } //only keep 2 decimal places of dollars

type database map[string]dollars //databaseis a map of items and their dollar values

func (db database) list(w http.ResponseWriter, req *http.Request) {
	mutex.RLock()                 //locks for reading
	for item, price := range db { //print items in the database
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	mutex.RUnlock()
}
func (db database) price(w http.ResponseWriter, req *http.Request) {
	mutex.RLock()
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404: item not in database
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	mutex.RUnlock()
}
func (db database) create(w http.ResponseWriter, req *http.Request) {
	mutex.Lock() //locks for writing
	item := req.URL.Query().Get("item")
	newPrice := req.URL.Query().Get("price")
	f, _ := strconv.ParseFloat(newPrice, 32)
	db[item] = dollars(f) //add new item to database
	fmt.Fprint(w, "New item created: ", item, ". It's price is now set to: ", f)
	mutex.Unlock()
}
func (db database) update(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	//change value of key
	item := req.URL.Query().Get("item")
	newPrice := req.URL.Query().Get("price")
	f, _ := strconv.ParseFloat(newPrice, 32)
	if price, ok := db[item]; ok {
		_ = price
		fmt.Fprint(w, "Price of ", item, " is now set to: ", f)
		db[item] = dollars(f)

	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "Error: No such item: %q\n", item)
	}
	mutex.Unlock()
}
func (db database) delete(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()

	itemToDelete := req.URL.Query().Get("item")

	delete(db, itemToDelete)

	fmt.Fprintf(w, itemToDelete, "has been deleted")

	mutex.Unlock()
}
