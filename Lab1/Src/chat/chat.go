package chat

import (
	"container/list"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

// timestamp varibale usada para timestamp al momento de escribir en ordenes.csv
var timestamp []time.Time

// data variable usada para escribir en el archivo ordenes.csv
var data []string

// numero variable para el seguimiento
var numero int = 1

// NumeroPaquete variable para id-paquete
var NumeroPaquete int = 1

// Server funcion mas diabla
type Server struct {
}

// Paquetes struct para guardar los paquetes en las colas
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

// PaqueteSeguimiento struct para guardar los paqutes en memoria
type PaqueteSeguimiento struct {
	IDPaquete   string
	Seguimiento string
	IDCamion    string
	Intentos    int
	Estado      string
	Tipo        string
	Valor       string
}

// PaqueteFinanza struct para pasar los datos a finanzas
type PaqueteFinanza struct {
	IDPaquete string
	Intentos  int
	Entrega   string
	Monto     int
}

// diccionario para guardar el id del paquetete asociado a un numero de seguimiento
var diccionario = make(map[string]int)

// Retail -> cola de paquetes que funciona como cola para paquetes tipo retail
var Retail = list.New()

// Prioritario -> cola de paquetes que funciona como cola para paquetes tipo retail
var Prioritario = list.New()

// Normal -> cola de paquetes que funciona como cola para paquetes tipo retail
var Normal = list.New()

// Memoria -> guarda la informacion de los paquetes para los seguimientos
var Memoria []PaqueteSeguimiento

// escribir escribe en el csv de ordenes todas las ordenes ingresada por los clientes

func escribir(archivo string, numeropaquete string, prioritario string, nombre string, valor string, origen string, destino string, numeroseguimiento string) {
	csvfile, err := os.OpenFile(archivo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	row := []string{timestamp, numeropaquete, prioritario, nombre, valor, origen, destino, numeroseguimiento}
	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Write(row)
	csvwriter.Flush()
	csvfile.Close()
}

// IngresarOrden funcion que recive una orden del cliente, escribe en el archivo y responde con el
// codigo de seguimiento de la orden solo si esta es de tipo pyme, ademas crea los paquetes en memoria y
// los manda a la cola correspondiente
func (s *Server) IngresarOrden(ctx context.Context, in *Message) (*MessageResponse, error) {

	mensaje := strings.Split(in.Mensaje, "#")
	_, nombre, valor, origen, destino, prioritario := mensaje[0], mensaje[1], mensaje[2], mensaje[3], mensaje[4], mensaje[5]

	NumeroSeguimiento := "N" + strconv.Itoa(numero)

	if prioritario == "0" {
		prioritario = "normal"
		numero++
	} else if prioritario == "1" {
		prioritario = "prioritario"
		numero++
	} else {
		prioritario = "retail"
		NumeroSeguimiento = "0"
	}

	paquete := Paquetes{
		strconv.Itoa(NumeroPaquete),
		NumeroSeguimiento,
		prioritario,
		valor,
		1,
		"En bodega",
		origen,
		destino}

	seg := PaqueteSeguimiento{
		strconv.Itoa(NumeroPaquete),
		NumeroSeguimiento,
		"0",
		1,
		"En bodega",
		prioritario,
		valor}

	// Escribir las ordenes de los clientes en el archivo ordenes.csv
	escribir("ordenes.csv", strconv.Itoa(NumeroPaquete), prioritario, nombre, valor, origen, destino, NumeroSeguimiento)

	// Poner paquetes en la cola correspondiente
	if prioritario == "normal" {

		Normal.PushBack(paquete)

	} else if prioritario == "prioritario" {

		Prioritario.PushBack(paquete)

	} else if prioritario == "retail" {

		Retail.PushBack(paquete)
	}

	Memoria = append(Memoria, seg)
	diccionario[NumeroSeguimiento] = NumeroPaquete

	NumeroPaquete++

	log.Printf("Orden recibida del cliente: %s", strings.Split(in.Mensaje, "#"))
	return &MessageResponse{Respuesta: NumeroSeguimiento}, nil
}

// SeguirPedido funcion que recibe un numero de seguimiento y retorna el estado del pedido
func (s *Server) SeguirPedido(ctx context.Context, in *Message) (*MessageResponse, error) {

	log.Printf("Orden de seguimiento del cliente: %s", in.Mensaje)

	num := diccionario[in.Mensaje]
	num--
	return &MessageResponse{Respuesta: Memoria[num].Estado}, nil
}

// VerificarPedido funcion que verifica si hay pedidos en la cola y responde con la cantidad
// de pedidos por cada cola (Cantidad de paquetes en cola retail, prioritario y normal)
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

// PedirPaquete funcion en la que un camion pide un paquete de alguno de los tres tipos
// y lo saca de la cola correspondiente enviando toda la informacion del paquete al camion
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

	// Si no hay paquetes entonces retorna ""
	return &Paquete{}, nil
}

// ActualizarPaquete funcion que actualiza el estado de un paquete y si los camiones vuelven a la central
// Logistica recibe la informacion y mediante rabbitmq y json envia la informacion necesaria a finanzas
func (s *Server) ActualizarPaquete(ctx context.Context, in *Message) (*MessageResponse, error) {
	log.Printf("Actualizar paquete: %s", in.Mensaje)

	mensaje := strings.Split(in.Mensaje, "#")
	id := mensaje[0]
	camion := mensaje[1]
	intentos := mensaje[2]
	estado := mensaje[3]

	num, err := strconv.Atoi(id)
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
		q, err := ch.QueueDeclare(
			"hello-queue", // name
			false,         // durable
			false,         // delete when unused
			false,         // exclusive
			false,         // no-wait
			nil,           // arguments
		)
		// We can print out the status of our Queue here
		// this will information like the amount of messages on
		// the queue

		// Handle any errors if we were unable to create the queue
		if err != nil {
			log.Fatalf("Handle any errors if we were unable to create the queue : %s", err)
		}

		IDPaque := Memoria[num].IDPaquete
		Inten := Memoria[num].Intentos
		Entreg := Memoria[num].Estado
		Tip := Memoria[num].Tipo
		Val := Memoria[num].Valor

		var Mont int

		Valo, err := strconv.Atoi(Val)
		if err != nil {
			log.Fatalf("Error al pasar el string a numero: %s", err)
		}

		if estado == "Entregado" {
			Mont = Valo
		} else if estado == "No entregado" {

			if Tip == "normal" {
				Mont = 0
			}

			if Tip == "prioritario" {
				Mont = int(30 * Valo / 100)
			}

			if Tip == "retail" {
				Mont = Valo
			}

		}

		mens := PaqueteFinanza{IDPaquete: IDPaque, Intentos: Inten, Entrega: Entreg, Monto: Mont}

		body, err := json.Marshal(mens)
		if err != nil {
			log.Fatalf("Error encoding JSON :%s", err)
		}
		// attempt to publish a message to the queue!
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})

		if err != nil {
			log.Fatalf("Error in publish: %s", err)
		}

	}

	return &MessageResponse{Respuesta: "Paquete actualizado"}, nil

}
