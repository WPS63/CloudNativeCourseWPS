package lru

import (
	"errors"
	"fmt"
)

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Put(interface{}, interface{}) error
}

type lruCache struct {
	size      int               //max size of cache
	remaining int               //remaining capacity
	cache     map[string]string //actual storage of data
	queue     []string          //for keeping tracks of least recently used
}

func NewCache(size int) Cacher {
	return &lruCache{size: size, remaining: size, cache: make(map[string]string), queue: make([]string, size)}
}

func (lru *lruCache) Get(key interface{}) (interface{}, error) {
	if _, ok := lru.cache[key.(string)]; ok {
		if lru.remaining == 0 {
			lru.queue[lru.size-1] = key.(string) //add to queue as highest available index
		} else {
			lru.queue[lru.size-lru.remaining] = key.(string) //add to queue as current index
		}
		return lru.cache[key.(string)], nil //return the value at the given key
	}
	return "-1", errors.New("That key is not in the map")
}

func (lru *lruCache) Put(key, val interface{}) error {
	if lru.remaining < 0 {
		return errors.New("Capacity error occurred, the cache is already full")
	}
	if lru.remaining == 0 {
		delete(lru.cache, lru.queue[0])   //delete the LRU from the cache
		lru.qDel(lru.queue[0])            //delete the LRU(head) from queue (which reduces queue slice size by one)
		lru.queue = append(lru.queue, "") //append empty string to queue to make it the original size again
		fmt.Print("Empty queue index amended: ")
		fmt.Println(lru.queue)
		lru.queue[lru.size-1] = key.(string) //now add the key to the tail of queue
		fmt.Print("Now the queue is: ")
		fmt.Println(lru.queue)
	} else {
		lru.queue[lru.size-lru.remaining] = key.(string) //if capacity isn't max, just add to slice
	}
	if lru.remaining > 0 { //prevents remaining from going below zero
		lru.remaining--
	}
	lru.cache[key.(string)] = val.(string) //insert into cache
	return nil
}

// Delete element from queue
func (lru *lruCache) qDel(ele string) {
	fmt.Print("The queue is full: ")
	fmt.Println(lru.queue)
	for i := 0; i < len(lru.queue); i++ {
		if lru.queue[i] == ele {
			oldlen := len(lru.queue)
			copy(lru.queue[i:], lru.queue[i+1:])
			lru.queue = lru.queue[:oldlen-1]
			break
		}
	}
	fmt.Print("The head is deleted: ")
	fmt.Println(lru.queue)
}
