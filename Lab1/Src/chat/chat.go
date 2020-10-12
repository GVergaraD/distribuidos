package chat

import (
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// Variable incremental para el numero de seguimiento
var numero int
var timestamp []time.Time

// Server funcion mas diabla
type Server struct {
}

// SayHello funcion que recive una orden del cliente, escribe en el archivo y responde el
// codigo de seguimiento de la orden
func (s *Server) SayHello(ctx context.Context, in *Message) (*MessageResponse, error) {
	// timestamp para meter al archivo
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("Receive message body from client: %s", strings.Split(in.Body, "-"))
	log.Printf(timestamp)
	numero++
	NumeroSeguimiento := strconv.Itoa(numero)
	return &MessageResponse{Seguimiento: NumeroSeguimiento}, nil
}
