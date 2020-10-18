package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/streadway/amqp"
)

// Total guarda el total
var Total int = 0

// Ganancias muestra las ganancias
var Ganancias int = 0

//Perdidas muestra las perdidas
var Perdidas int = 0

// InfoPaquete para pasar a finanzas
type InfoPaquete struct {
	IDPaquete string
	Intentos  int
	Entrega   string
	Monto     int
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func guardar(archivo string, idpaquete string, intentos string, entrega string, monto string) {
	csvfile, err := os.OpenFile(archivo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	row := []string{idpaquete, intentos, entrega, monto}
	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Write(row)
	csvwriter.Flush()
	csvfile.Close()
}

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	csvfile, err := os.OpenFile("balance.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	row := []string{"IDPaquete", "Intentos", "Entrega", "Ganancias/Perdidas"}
	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Write(row)
	csvwriter.Flush()
	csvfile.Close()

	q, err := ch.QueueDeclare(
		"hello-queue", // name
		false,         // durable
		false,         // delete when usused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			Info := &InfoPaquete{}
			err := json.Unmarshal(d.Body, Info)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
			IDPaquete := Info.IDPaquete
			Intentos := Info.Intentos
			Entrega := Info.Entrega
			Monto := Info.Monto

			if Monto < 0 {
				Perdidas = Perdidas + Monto
			} else if Monto > 0 {
				Ganancias = Ganancias + Monto
			}

			Total = Ganancias - Perdidas

			StringIntentos := strconv.Itoa(Intentos)
			StringMonto := strconv.Itoa(Monto)

			guardar("balance.csv", IDPaquete, StringIntentos, Entrega, StringMonto)
			log.Printf("Total: %s | Ganancias %s | Perdidas: %s", strconv.Itoa(Total), strconv.Itoa(Ganancias), strconv.Itoa(Perdidas))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
