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
	http.HandleFunc("/delete", db.delete) //additional handlers for new functions
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float64             //declare type dollars
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) } //only keep 2 decimal places of dollars
type database map[string]dollars //database is a map of items and their dollar values
var mutex = sync.RWMutex{}       //declare mutex to use for locks

//curl "http://localhost:8000/list"
func (db database) list(w http.ResponseWriter, req *http.Request) {
	mutex.RLock()                 //locks for reading
	for item, price := range db { //print items in the database
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	mutex.RUnlock()
}

//curl "http://localhost:8000/price?item=socks"
func (db database) price(w http.ResponseWriter, req *http.Request) {
	mutex.RLock()
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) //item not in database
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	mutex.RUnlock()
}

//curl "http://localhost:8000/create?item=pants&price=20"
func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Price: %v\n", err)
		return
	} else if p <= 0 {
		fmt.Fprintf(w, "That price is too low. Try again.")
		return
	}
	if _, found := db[item]; found {
		fmt.Fprint(w, "That item is already in the db so I updated it's value for you.\n")
		mutex.Lock() //locks for writing
		db[item] = dollars(p)
		mutex.Unlock() //locks for writing
		fmt.Fprint(w, "Updated item: ", item, ". Item price:  ", p, " \n")
	} else {
		mutex.Lock() //locks for writing
		db[item] = dollars(p)
		mutex.Unlock() //locks for writing
		fmt.Fprint(w, "Created item: ", item, ". Item price:  ", p, " \n")
	}
}

//curl "http://localhost:8000/update?item=socks&price=7"
func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Price: %v\n", err)
		return
	} else if p <= 0 {
		fmt.Fprintf(w, "That price is too low. Try again.")
		return
	}
	if _, found := db[item]; found {
		mutex.Lock()
		db[item] = dollars(p) //convert float to type dollars and store in map
		mutex.Unlock()
		fmt.Fprint(w, "Updated item: ", item, ". Item price:  ", p, " \n")
	} else {
		w.WriteHeader(http.StatusNotFound) // 404: item not in database
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

//curl "http://localhost:8000/delete?item=socks"
func (db database) delete(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	itemToDelete := req.URL.Query().Get("item")
	if _, found := db[itemToDelete]; found {

		delete(db, itemToDelete)

		fmt.Fprintf(w, "Deleted item %s\n", itemToDelete)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404: item not in database
		fmt.Fprintf(w, "no such item: %q\n", itemToDelete)
	}
	mutex.Unlock()
}
