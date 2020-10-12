package chat

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// variable declaration using a type
var timestamp []time.Time
var data []string

// numero variable para el seguimiento
var numero int = 1

// NumeroPaquete variable para id-paquete
var NumeroPaquete int = 100

// Server funcion mas diabla
type Server struct {
}

// SayHello funcion que recive una orden del cliente, escribe en el archivo y responde el
// codigo de seguimiento de la orden
func (s *Server) SayHello(ctx context.Context, in *Message) (*MessageResponse, error) {

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	mensaje := strings.Split(in.Body, "#")
	_, nombre, valor, origen, destino, prioritario := mensaje[0], mensaje[1], mensaje[2], mensaje[3], mensaje[4], mensaje[5]

	if prioritario == "0" {
		prioritario = "normal"
	} else if prioritario == "1" {
		prioritario = "prioritario"
	} else {
		prioritario = "retail"
	}

	NumeroSeguimiento := strconv.Itoa(numero)
	row := []string{timestamp, "N" + strconv.Itoa(NumeroPaquete), prioritario, nombre, valor, origen, destino, NumeroSeguimiento}
	numero++
	NumeroPaquete++

	csvfile, err := os.OpenFile("ordenes.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Write(row)
	csvwriter.Flush()
	csvfile.Close()

	log.Printf("Orden recibida del cliente: %s", strings.Split(in.Body, "#"))
	return &MessageResponse{Seguimiento: NumeroSeguimiento}, nil
}
