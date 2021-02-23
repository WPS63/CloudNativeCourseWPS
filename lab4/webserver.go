package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := database{"shoes": 50, "socks": 5} //create an instance of type database
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/delete", db.delete)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float32 //declare type dollars

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) } //only keep 2 decimal places of dollars

type database map[string]dollars //databases are maps of items and their dollar value

func (db database) list(w http.ResponseWriter, req *http.Request) { //prints all items in the database
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}
func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
func (db database) create(w http.ResponseWriter, req *http.Request) {
	//add new item to
}
func (db database) update(w http.ResponseWriter, req *http.Request) {

}
func (db database) delete(w http.ResponseWriter, req *http.Request) {
	itemToDelete := req.URL.Query().Get("item")

	for item, price := range db { //this is here to show list before delete(remove before demo)
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}

	delete(db, itemToDelete)

	for item, price := range db { //this is here to show list after delete(remove before demo)
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}
