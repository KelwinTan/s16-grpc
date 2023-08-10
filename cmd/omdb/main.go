package main

import (
	"log"
	"net"

	grpcS16 "github.com/KelwinTan/s16-grpc/api/grpc/omdb"
	"github.com/KelwinTan/s16-grpc/api/proto/v1/omdb"
	"github.com/KelwinTan/s16-grpc/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	PORT = ":50051"
)

func main() {

	app.Init()

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	omdb.RegisterOMDBServiceServer(s, &grpcS16.Server{})

	log.Printf("Server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
