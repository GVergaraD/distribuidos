package main

import (
	"fmt"
	"log"
	"net"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9003))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	lis2, err := net.Listen("tcp", fmt.Sprintf(":%d", 9003)) //ip nodo3
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := chat.Server{}

	grpcServer := grpc.NewServer()

	chat.RegisterLibroServiceServer(grpcServer, &s)
	chat.RegisterChunkServiceServer(grpcServer, &s)
	chat.RegisterPropuestaServiceServer(grpcServer, &s)
	chat.RegisterContactarServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	if err := grpcServer.Serve(lis2); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
