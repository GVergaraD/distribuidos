package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
	"google.golang.org/grpc"
)

func main() {

	// Escribir encabezados en ordenes.csv
	row := []string{"timestamp", "id-paquete", "tipo", "nombre", "valor", "origen", "destino", "seguimiento"}
	csvfile, err := os.OpenFile("ordenes.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Write(row)
	csvwriter.Flush()
	csvfile.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := chat.Server{}

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)
	chat.RegisterSeguimientoServer(grpcServer, &s)
	chat.RegisterVerificarServer(grpcServer, &s)
	chat.RegisterPedirServer(grpcServer, &s)
	chat.RegisterActualizarServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
