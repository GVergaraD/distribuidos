package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	// Abrir archivo
	csvfile, err := os.Open("pymes.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	csvReader := csv.NewReader(csvfile)

	// La primera linea tiene los encabezados
	if _, err := csvReader.Read(); err != nil {
		log.Fatalln("Error al leer la primera linea", err)
	}

	for {
		post, err := csvReader.Read()
		if err != nil {
			log.Println(err)
			// Will break on EOF.
			break
		}
		id, producto, valor, tienda, destino, prioritario := post[0], post[1], post[2], post[3], post[4], post[5]
		response, err := c.SayHello(context.Background(), &chat.Message{Body: id + "-" + producto + "-" + valor + "-" + tienda + "-" + destino + "-" + prioritario})
		if err != nil {
			log.Fatalf("Error al enviar la orden: %s", err)
		}
		log.Printf("Orden ingresada, numero de seguimiento: %s", response.Seguimiento)
		time.Sleep(2 * time.Second)
	}
}
