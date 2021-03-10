// Package main implements a server for movieinfo service.
package main

import (
	"context"
	"errors"
	"fmt"
	"labs/lab5/movieapi"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement movieapi.MovieInfoServer
type server struct {
	movieapi.UnimplementedMovieInfoServer
}

// Map representing a database
var moviedb = map[string][]string{"Pulp fiction": []string{"1994", "Quentin Tarantino", "John Travolta,Samuel Jackson,Uma Thurman,Bruce Willis"}}

// GetMovieInfo implements movieapi.MovieInfoServer
func (s *server) GetMovieInfo(ctx context.Context, in *movieapi.MovieRequest) (*movieapi.MovieReply, error) {
	title := in.GetTitle()
	log.Printf("Received: %v", title)
	reply := &movieapi.MovieReply{}
	if val, ok := moviedb[title]; !ok { // Title not present in database
		return reply, nil
	} else {
		if year, err := strconv.Atoi(val[0]); err != nil {
			reply.Year = -1
		} else {
			reply.Year = int32(year)
		}
		reply.Director = val[1]
		cast := strings.Split(val[2], ",")
		reply.Cast = append(reply.Cast, cast...)
	}
	return reply, nil
}
func (s *server) SetMovieInfo(ctx context.Context, message *movieapi.MovieData) (*movieapi.Status, error) {
	title := message.GetTitle() //these functions are declared automatically in the movieapi.pb.go
	year := message.GetYear()
	director := message.GetDirector()
	cast := message.GetCast()
	status := &movieapi.Status{}     //create a new status message to be returned to the client after
	yearS := fmt.Sprint(year)        //convert year into a string so you can add it to the string array
	castS := strings.Join(cast, ",") //combine the strings of cast slice into one string
	mapValue := make([]string, 0)
	mapValue = append(mapValue, yearS, director, castS) //append movie data to the slice
	if _, presentInDb := moviedb[title]; !presentInDb {
		moviedb[title] = mapValue //add slice to map
		status.Code = "success"
		return status, nil
	}
	status.Code = "failure"
	return status, errors.New("Movie already in db")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	movieapi.RegisterMovieInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
