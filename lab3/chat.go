// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct {
	//Converting type client
	clientChan chan<- string //an outgoing message channel for each client
	clientName string        //name of client
}

var (
	entering = make(chan client) //channel of client channels
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000") //listen on host port
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster() //sends the messages
	for {
		conn, err := listener.Accept() //new connection
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn) //create a channel for new connection
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.clientChan <- msg //send message to each cli in clients map
			}

		case cli := <-entering: //if the client enters then...
			clients[cli] = true //adds value: true to key: cli in client map
			cli.clientChan <- "Current chatters: "
			for c := range clients { //Displays current list of all clients in client map to new entry
				cli.clientChan <- c.clientName
			}

		case cli := <-leaving: //if the client leaves then...
			delete(clients, cli)  //delete from map
			close(cli.clientChan) //close channel
		}
	}
}

func handleConn(conn net.Conn) { //creates a channel for the client to receive messages from other clients

	var cli client            //cli is a new client instance
	ch := make(chan string)   //make a channel for outgoing client messages
	go clientWriter(conn, ch) //call clientWriter subroutine for current connection and newly made channel

	//Set the username for cli
	inputName := bufio.NewReader(conn)
	fmt.Fprintln(conn, "Enter name: ")
	name, _ := inputName.ReadString('\n')
	who := strings.TrimSuffix(name, "\n")

	cli.clientChan = ch              //cli's channel is 'ch'
	cli.clientName = who             //cli's name is 'who'
	ch <- "You are " + who           //tell cli their name
	messages <- who + " has arrived" //tell everyone else cli has arrived
	entering <- cli                  //send cli to 'entering' channel

	input := bufio.NewScanner(conn)
	for input.Scan() { //scan inputs until ctrl C is hit
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- cli                //send cli to 'leaving' channel
	messages <- who + " has left" //tell everyone that cli has left
	conn.Close()                  //close connection
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
