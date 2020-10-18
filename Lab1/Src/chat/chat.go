package chat

import (
	"container/list"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
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

// Paquetes -> Estructura de datos para el paquete
type Paquetes struct {
	IDPaquete   string
	Seguimiento string
	Tipo        string
	Valor       string
	Intentos    int
	Estado      string
	Origen      string
	Destino     string
}

// PaqueteSeguimiento funcion
type PaqueteSeguimiento struct {
	IDPaquete   string
	Seguimiento string
	IDCamion    string
	Intentos    int
	Estado      string
}

// Retail -> cola de paquetes que funciona como cola para paquetes tipo retail
var Retail = list.New()

// Prioritario -> cola de paquetes que funciona como cola para paquetes tipo retail
var Prioritario = list.New()

// Normal -> cola de paquetes que funciona como cola para paquetes tipo retail
var Normal = list.New()

// Memoria -> guarda la informacion de los paquetes para los seguimientos
var Memoria []PaqueteSeguimiento

// IngresarOrden funcion que recive una orden del cliente, escribe en el archivo y responde el
// codigo de seguimiento de la orden, ademas crea los paquetes y los manda a la cola correspondiente
func (s *Server) IngresarOrden(ctx context.Context, in *Message) (*MessageResponse, error) {

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	mensaje := strings.Split(in.Mensaje, "#")
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

	csvfile, err := os.OpenFile("ordenes.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Write(row)
	csvwriter.Flush()
	csvfile.Close()

	paquete := Paquetes{
		"N" + strconv.Itoa(NumeroPaquete),
		NumeroSeguimiento,
		prioritario,
		valor,
		0,
		"En bodega",
		origen,
		destino}

	if prioritario == "normal" {

		Normal.PushBack(paquete)

	} else if prioritario == "prioritario" {

		Prioritario.PushBack(paquete)

	} else if prioritario == "retail" {

		Retail.PushBack(paquete)
	}

	seg := PaqueteSeguimiento{
		"N" + strconv.Itoa(NumeroPaquete),
		NumeroSeguimiento,
		"0",
		1,
		"En bodega"}

	Memoria = append(Memoria, seg)

	numero++
	NumeroPaquete++

	log.Printf("Orden recibida del cliente: %s", strings.Split(in.Mensaje, "#"))
	return &MessageResponse{Respuesta: NumeroSeguimiento}, nil
}

// SeguirPedido funcion
func (s *Server) SeguirPedido(ctx context.Context, in *Message) (*MessageResponse, error) {
	log.Printf("Orden de seguimiento del cliente: %s", in.Mensaje)

	num, err := strconv.Atoi(in.Mensaje)
	if err != nil {
		log.Fatalf("Error al pasar el string a numero: %s", err)
	}
	num--
	return &MessageResponse{Respuesta: Memoria[num].Estado}, nil
}

// VerificarPedido funcion que verifica si hay pedidos en la cola
func (s *Server) VerificarPedido(ctx context.Context, in *Message) (*Colas, error) {

	log.Printf("El camion pregunta: %s", in.Mensaje)

	if Retail.Len() >= 1 || Prioritario.Len() >= 1 || Normal.Len() >= 1 {

		return &Colas{Respuesta: "Si",
			Retail:      strconv.Itoa(Retail.Len()),
			Prioritario: strconv.Itoa(Prioritario.Len()),
			Normal:      strconv.Itoa(Normal.Len())}, nil

	}

	return &Colas{Respuesta: "No"}, nil
}

// PedirPaquete funcion que pide un paquete y lo saca de la cola
func (s *Server) PedirPaquete(ctx context.Context, in *Message) (*Paquete, error) {

	log.Printf("El camion pregunta: %s", in.Mensaje)

	if in.Mensaje == "Retail" {
		front := Retail.Front()
		elemento := Paquetes(front.Value.(Paquetes))
		idpaquete := elemento.IDPaquete
		tipo := elemento.Tipo
		valor := elemento.Valor
		intentos := strconv.Itoa(elemento.Intentos)
		origen := elemento.Origen
		destino := elemento.Destino
		seguimiento := elemento.Seguimiento
		Retail.Remove(front)

		return &Paquete{ID: idpaquete,
			Tipo:        tipo,
			Valor:       valor,
			Origen:      origen,
			Destino:     destino,
			Intentos:    intentos,
			Seguimiento: seguimiento}, nil
	}

	if in.Mensaje == "Prioritario" {
		front := Prioritario.Front()
		elemento := Paquetes(front.Value.(Paquetes))
		idpaquete := elemento.IDPaquete
		tipo := elemento.Tipo
		valor := elemento.Valor
		intentos := strconv.Itoa(elemento.Intentos)
		origen := elemento.Origen
		destino := elemento.Destino
		seguimiento := elemento.Seguimiento
		Prioritario.Remove(front)

		return &Paquete{ID: idpaquete,
			Tipo:        tipo,
			Valor:       valor,
			Origen:      origen,
			Destino:     destino,
			Intentos:    intentos,
			Seguimiento: seguimiento}, nil
	}

	if in.Mensaje == "Normal" {
		front := Normal.Front()
		elemento := Paquetes(front.Value.(Paquetes))
		idpaquete := elemento.IDPaquete
		tipo := elemento.Tipo
		valor := elemento.Valor
		intentos := strconv.Itoa(elemento.Intentos)
		origen := elemento.Origen
		destino := elemento.Destino
		seguimiento := elemento.Seguimiento
		Normal.Remove(front)

		return &Paquete{ID: idpaquete,
			Tipo:        tipo,
			Valor:       valor,
			Origen:      origen,
			Destino:     destino,
			Intentos:    intentos,
			Seguimiento: seguimiento}, nil
	}

	return &Paquete{}, nil
}

// ActualizarPaquete funcion
func (s *Server) ActualizarPaquete(ctx context.Context, in *Message) (*MessageResponse, error) {
	log.Printf("Actualizar paquete: %s", in.Mensaje)

	// in.Mensaje = "1#camion1#0#En camino"
	// in.Mensaje = "1#camion1#3#Entrega"
	// in.Mensaje = seguimiento#camion#intentos#estado
	mensaje := strings.Split(in.Mensaje, "#")
	seguimiento := mensaje[0]
	camion := mensaje[1]
	intentos := mensaje[2]
	estado := mensaje[3]

	num, err := strconv.Atoi(seguimiento)
	if err != nil {
		log.Fatalf("Error al pasar el string a numero: %s", err)
	}
	num--

	NIntentos, err := strconv.Atoi(intentos)
	if err != nil {
		log.Fatalf("Error al pasar el string a numero: %s", err)
	}

	if camion != "0" {
		Memoria[num].IDCamion = camion
	}

	if intentos != "0" {
		Memoria[num].Intentos = NIntentos
	}

	if estado != "0" {
		Memoria[num].Estado = estado
	}

	if estado == "Entregado" || estado == "No entregado" {

		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Fatalf("Failed Initializing Broker Connection : %s", err)
		}

		// Let's start by opening a channel to our RabbitMQ instance
		// over the connection we have already established
		ch, err := conn.Channel()
		if err != nil {

			log.Fatalf("Failed Initializing Connection : %s", err)
		}
		defer ch.Close()

		// with this channel open, we can then start to interact
		// with the instance and declare Queues that we can publish and
		// subscribe to
		_, err = ch.QueueDeclare(
			"TestQueue",
			false,
			false,
			false,
			false,
			nil,
		)
		// We can print out the status of our Queue here
		// this will information like the amount of messages on
		// the queue

		// Handle any errors if we were unable to create the queue
		if err != nil {
			log.Fatalf("Handle any errors if we were unable to create the queue : %s", err)
		}

		// attempt to publish a message to the queue!
		err = ch.Publish(
			"",
			"TestQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("Hello World"),
			},
		)

		if err != nil {
			log.Fatalf("Error in publish: %s", err)
		}

	}

	return &MessageResponse{Respuesta: "Paquete actualizado : %s"}, nil

}
