// Package main imlements a client for movieinfo service
package main

import (
	"context"
	"log"
	"os"
	"time"

	"labs/lab5/movieapi"

	"google.golang.org/grpc"
)

const (
	address            = "localhost:50051"
	defaultTitle       = "Pulp fiction"
	newTitle           = "Lion King"
	newYear      int32 = 1994
	newDirector        = "Rob Minkoff"
)

var newCast = []string{"James Earl Jones", "Jeremy Irons"}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Contact the server and print out its response.
	title := defaultTitle
	if len(os.Args) > 1 {
		title = os.Args[1]
	}
	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s %d %s %v", title, r.GetYear(), r.GetDirector(), r.GetCast())

	w, err := c.SetMovieInfo(ctx, &movieapi.MovieData{Title: newTitle, Year: newYear, Director: newDirector, Cast: newCast})

	log.Print("Status of set is:" + w.GetCode())

	q, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: newTitle})

	log.Printf("Movie Info for %s %d %s %v", newTitle, q.GetYear(), q.GetDirector(), q.GetCast())
}
