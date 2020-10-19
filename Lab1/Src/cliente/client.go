package main

import (
	"encoding/csv"
	"fmt"
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
	s := chat.NewSeguimientoClient(conn)

	var value int

	for {
		fmt.Println("\nIgrese opcion a realizar")
		fmt.Println("1.- Ingresar pedido")
		fmt.Println("2.- Seguir pedido")
		fmt.Printf("3.- Salir\n\n")
		fmt.Scanln(&value)

		if value == 1 {

			fmt.Println("\nComportamiento a seguir (pymes/retail): ")
			var comportamiento string
			fmt.Scanln(&comportamiento)

			fmt.Println("\nTiempo de espera en segundos entre ordenes: ")
			var tiempo int
			fmt.Scanln(&tiempo)

			// Abrir archivo
			csvfile, err := os.Open(comportamiento + ".csv")
			if err != nil {
				log.Fatalln("Couldn't open the csv file", err)
			}

			csvReader := csv.NewReader(csvfile)

			// La primera linea tiene los encabezados
			if _, err := csvReader.Read(); err != nil {
				log.Fatalln("Error al leer la primera linea", err)
			}

			if comportamiento == "pymes" {
				for {
					post, err := csvReader.Read()
					if err != nil {
						log.Println(err)
						// Will break on EOF.
						break
					}

					id, producto, valor, tienda, destino, prioritario := post[0], post[1], post[2], post[3], post[4], post[5]
					response, err := c.IngresarOrden(context.Background(), &chat.Message{Mensaje: id + "#" + producto + "#" + valor + "#" + tienda + "#" + destino + "#" + prioritario})
					if err != nil {
						log.Fatalf("Error al enviar la orden: %s", err)
					}
					log.Printf("Orden ingresada, numero de seguimiento: %s", response.Respuesta)
					time.Sleep(time.Duration(tiempo) * time.Second)
				}
			} else {
				for {
					post, err := csvReader.Read()
					if err != nil {
						log.Println(err)
						// Will break on EOF.
						break
					}

					id, producto, valor, tienda, destino, prioritario := post[0], post[1], post[2], post[3], post[4], "-1"
					_, err = c.IngresarOrden(context.Background(), &chat.Message{Mensaje: id + "#" + producto + "#" + valor + "#" + tienda + "#" + destino + "#" + prioritario})
					if err != nil {
						log.Fatalf("Error al enviar la orden: %s", err)
					}
					log.Printf("Orden ingresada")
					time.Sleep(time.Duration(tiempo) * time.Second)
				}
			}
		} else if value == 2 {

			fmt.Println("\nIngrese numero de seguimiento: ")
			var seguimiento string
			fmt.Scanln(&seguimiento)

			response, err := s.SeguirPedido(context.Background(), &chat.Message{Mensaje: seguimiento})
			if err != nil {
				log.Fatalf("Error al seguir pedido: %s", err)
			}
			log.Printf("Estado de su pedido: %s", response.Respuesta)

		} else if value == 3 {
			break
		}
	}
}
