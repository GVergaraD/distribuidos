package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
)

// Variables para gestionar los paquetes en los camiones
var retail int
var prioritario int
var normal int
var p11 int = 0
var p12 int = 0
var p21 int = 0
var p22 int = 0
var p31 int = 0
var p32 int = 0
var tipo11 string
var tipo12 string
var tipo21 string
var tipo22 string
var tipo31 string
var tipo32 string

func savearchivo(camion string, idpaquete string, tipo string, valor string, origen string, destino string, intentos int, fallo int) {
	csvfile, err := os.OpenFile(camion, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	csvwriter := csv.NewWriter(csvfile)
	if fallo == 0 {
		row := []string{idpaquete, tipo, valor, origen, destino, strconv.Itoa(intentos), timestamp}
		csvwriter.Write(row)
	} else {
		row := []string{idpaquete, tipo, valor, origen, destino, strconv.Itoa(intentos), "0"}
		csvwriter.Write(row)
	}

	csvwriter.Flush()
	csvfile.Close()
}

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewVerificarClient(conn)
	p := chat.NewPedirClient(conn)
	a := chat.NewActualizarClient(conn)

	row := []string{"id-paquete", "tipo", "valor", "origen", "destino", "intentos", "fecha-entrega"}

	archivos := []string{"camion1.csv", "camion2.csv", "camion3.csv"}

	for i := 0; i < 3; i++ {

		csvfile, err := os.OpenFile(archivos[i], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		csvwriter := csv.NewWriter(csvfile)
		csvwriter.Write(row)
		csvwriter.Flush()
		csvfile.Close()
	}

	//ejemplo de pedir un paquete
	//respuesta, err := p.PedirPaquete(context.Background(), &chat.Message{Mensaje: "Retail"})
	//if err != nil {
	//	log.Fatalf("Error al verificar si hay pedidos: %s", err)
	//}
	//log.Printf("Respuesta de logistica: %s", respuesta.ID)
	//log.Printf("Respuesta de logistica: %s", response.ID, response.Tipo, response.Valor, response.Origen, response.Destino, response.Intentos (este es un string por el momento))
	// Estado inicial, los 3 camiones se encuentran disponibles

	var ch chan string = make(chan string)
	var ch2 chan string = make(chan string)
	var ch3 chan string = make(chan string)
	var chdisp1 chan string = make(chan string)
	var chdisp2 chan string = make(chan string)
	var chdisp3 chan string = make(chan string)
	var chp chan string = make(chan string)
	var chp2 chan string = make(chan string)
	var chp3 chan string = make(chan string)
	var espera int = 3  //tiempo de espera por 2do paquete
	var entrega int = 1 //tiempo de entrega
	fmt.Println("\ntiempo de espera por 2do paquete ( mayor a 0 ): ")
	fmt.Scanln(&espera)
	fmt.Println("\ntiempo de entrega: ")
	fmt.Scanln(&entrega)

	//defer close(ch2)
	go func() { //camion 1
		capacidad1 := 2
		random := 0
		intentos := 1
		intentos2 := 1
		var id1 string
		var tipo1 string
		var valor1 string
		var origen1 string
		var destino1 string
		var id2 string
		var tipo2 string
		var valor2 string
		var origen2 string
		var destino2 string
		//help=1
		//log.Printf("Camion 1 listo")
		for {
			//time.Sleep(200 * time.Millisecond)
			if capacidad1 == 2 {
				chdisp1 <- "ready"
			}
			select {
			case <-ch:
				//procede a rutina
				//fmt.Println("Camion 1 en la cola")
				//se recibe un paquete
				if <-chdisp1 == "yes" {
					capacidad1--
					//fmt.Println("Camion 1 espera 2do paquete")
				}
				//pedir
				select {
				case <-chp:
					respuesta1, err := p.PedirPaquete(context.Background(), &chat.Message{Mensaje: tipo11})
					if err != nil {
						log.Fatalf("Error al pedir paquete: %s", err)
					}
					log.Printf("Respuesta de logistica: %s", respuesta1.ID)
					p11, err = strconv.Atoi(respuesta1.Valor)
					id1 = respuesta1.ID
					tipo1 = respuesta1.Tipo
					valor1 = respuesta1.Valor
					origen1 = respuesta1.Origen
					destino1 = respuesta1.Destino
					_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(0) + "#En camino"})
					if err != nil {
						log.Fatalf("Error al actualizar paquete: %s", err)
					}
				}
				//time.Sleep(200 * time.Millisecond)
				select {
				case <-chdisp1:
					capacidad1--
					//pedir2
					select {
					case <-chp:
						if retail >= 1 || prioritario >= 1 || normal >= 1 {
							//fmt.Println("pidiendo 2do paquete recibido;", tipo12)
							respuesta2, err := p.PedirPaquete(context.Background(), &chat.Message{Mensaje: tipo12})
							if err != nil {
								log.Fatalf("Error al pedir paquete: %s", err)
							}
							log.Printf("Respuesta de logistica: %s", respuesta2.ID)
							p12, err = strconv.Atoi(respuesta2.Valor)
							id2 = respuesta2.ID
							tipo2 = respuesta2.Tipo
							valor2 = respuesta2.Valor
							origen2 = respuesta2.Origen
							destino2 = respuesta2.Destino
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(0) + "#En camino"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							//fmt.Println("2do paquete recibido")
						}
					}
				case <-time.After(time.Duration(espera) * time.Second): //tiempo de espera por 2do paquete
					//fmt.Println("no llegó 2do pedido")
				}
				//fmt.Println("Camion 1 se va")
				//time.Sleep(time.Second * 20)
				//funcion para entregar paquetes
				//se debe ver cual genera mas ingreso
				switch {
				case p12 == 0: //lego solo 1 paquete
					switch {
					case tipo11 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p11 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos++
							//fmt.Println("Reintento")
						}
						if p11 != 0 {
							savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p11 = 0
						}
					case tipo11 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p11 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p11-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p11-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p11 = 0
								break
							}
						}
						if p11 != 0 {
							savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p11 = 0
						}
					}
				case p11 > p12 && p12 != 0: //el 1er paquete vale más
					switch {
					case tipo11 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p11 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos++
							//fmt.Println("Reintento")
						}
						if p11 != 0 {
							savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p11 = 0
						}
					case tipo11 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p11 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p11-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p11-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p11 = 0
								break
							}
						}
						if p11 != 0 {
							savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p11 = 0
						}
					}
					switch {
					case tipo12 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p12 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos2++
							//fmt.Println("Reintento")
						}
						if p12 != 0 {
							savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p12 = 0
						}
					case tipo12 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p12 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p12-10 > 0:
								intentos2++
								//fmt.Println("Reintento")
							case p12-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p12 = 0
								break
							}
						}
						if p12 != 0 {
							savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p12 = 0
						}
					}
				case p12 > p11: //el 2do paquete vale mas
					switch {
					case tipo12 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p12 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos2++
							//fmt.Println("Reintento")
						}
						if p12 != 0 {
							savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p12 = 0
						}
					case tipo12 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p12 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p12-10 > 0:
								intentos2++
								//fmt.Println("Reintento")
							case p12-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p12 = 0
								break
							}
						}
						if p12 != 0 {
							savearchivo("camion1.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion1#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p12 = 0
						}
					}
					switch {
					case tipo11 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p11 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos++
							//fmt.Println("Reintento")
						}
						if p11 != 0 {
							savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p11 = 0
						}
					case tipo11 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION
								savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p11 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p11-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p11-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p11 = 0
								break
							}
						}
						if p11 != 0 {
							savearchivo("camion1.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion1#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p11 = 0
						}
					}
				}
				//fmt.Println("Camion 1, Paquete:entregado")
				////fmt.Println("Paquete en la cola")
			}
			capacidad1 = 2
			intentos = 1
			intentos2 = 1
		}

	}()
	go func() { //camion 2
		capacidad2 := 2
		random := 0
		intentos := 1
		intentos2 := 1
		var id1 string
		var tipo1 string
		var valor1 string
		var origen1 string
		var destino1 string
		var id2 string
		var tipo2 string
		var valor2 string
		var origen2 string
		var destino2 string
		//log.Printf("Camion 2 listo")
		for {
			//time.Sleep(200 * time.Millisecond)
			if capacidad2 == 2 {
				chdisp2 <- "ready"
			}

			select {
			case <-ch2:

				//time.Sleep(time.Second * 2)
				//fmt.Println("Camion 2 en la cola")
				//se recibe un paquete
				if <-chdisp2 == "yes" {
					capacidad2--
					//fmt.Println("Camion 2 espera 2do paquete")
				}
				//pedir
				select {
				case <-chp2:
					respuesta1, err := p.PedirPaquete(context.Background(), &chat.Message{Mensaje: tipo21})
					if err != nil {
						log.Fatalf("Error al pedir paquete: %s", err)
					}
					log.Printf("Respuesta de logistica: %s", respuesta1.ID)
					p21, err = strconv.Atoi(respuesta1.Valor)
					id1 = respuesta1.ID
					tipo1 = respuesta1.Tipo
					valor1 = respuesta1.Valor
					origen1 = respuesta1.Origen
					destino1 = respuesta1.Destino
					_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(0) + "#En camino"})
					if err != nil {
						log.Fatalf("Error al actualizar paquete: %s", err)
					}
				}
				//time.Sleep(200 * time.Millisecond)
				select {
				case <-chdisp2:
					capacidad2--
					//pedir2
					select {
					case <-chp2:
						//fmt.Println("pidiendo 2do paquete recibido;", tipo22)
						if retail >= 1 || prioritario >= 1 || normal >= 1 {
							respuesta2, err := p.PedirPaquete(context.Background(), &chat.Message{Mensaje: tipo22})
							if err != nil {
								log.Fatalf("Error al pedir paquete: %s", err)
							}
							log.Printf("Respuesta de logistica: %s", respuesta2.ID)
							p22, err = strconv.Atoi(respuesta2.Valor)
							id2 = respuesta2.ID
							tipo2 = respuesta2.Tipo
							valor2 = respuesta2.Valor
							origen2 = respuesta2.Origen
							destino2 = respuesta2.Destino
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(0) + "#En camino"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							//fmt.Println("2do paquete recibido")
						}
					}
				case <-time.After(time.Duration(espera) * time.Second): //tiempo de espera por 2do paquete
					//fmt.Println("no llegó 2do pedido")
				}
				//fmt.Println("Camion 2 se va")
				//time.Sleep(time.Second * 20)
				//funcion para entregar paquetes
				//se debe ver cual genera mas ingreso
				switch {
				case p22 == 0: //lego solo 1 paquete
					switch {
					case tipo21 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p21 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos++
							//fmt.Println("Reintento")
						}
						if p21 != 0 {
							savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p21 = 0
						}
					case tipo21 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p21 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p21-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p21-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p21 = 0
								break
							}
						}
						if p21 != 0 {
							savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p21 = 0
						}
					}
				case p21 > p22 && p22 != 0: //el 1er paquete vale más
					switch {
					case tipo21 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p21 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos++
							//fmt.Println("Reintento")
						}
						if p21 != 0 {
							savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p21 = 0
						}
					case tipo21 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p21 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p21-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p21-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p21 = 0
								break
							}
						}
						if p21 != 0 {
							savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p21 = 0
						}
					}
					switch {
					case tipo22 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p22 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos2++
							//fmt.Println("Reintento")
						}
						if p22 != 0 {
							savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p22 = 0
						}
					case tipo22 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p22 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p22-10 > 0:
								intentos2++
								//fmt.Println("Reintento")
							case p22-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p22 = 0
								break
							}
						}
						if p22 != 0 {
							savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p22 = 0
						}
					}
				case p22 > p21: //el 2do paquete vale mas
					switch {
					case tipo22 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p22 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos2++
							//fmt.Println("Reintento")
						}
						if p22 != 0 {
							savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p22 = 0
						}
					case tipo22 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p22 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p22-10 > 0:
								intentos2++
								//fmt.Println("Reintento")
							case p22-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p22 = 0
								break
							}
						}
						if p22 != 0 {
							savearchivo("camion2.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p22 = 0
						}
					}
					switch {
					case tipo21 == "Retail":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p21 = 0
								break
							}
							//fmt.Println("intento fallido")
							intentos++
							//fmt.Println("Reintento")
						}
						if p21 != 0 {
							savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p21 = 0
						}
					case tipo21 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION2
								savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p21 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p21-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p21-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p21 = 0
								break
							}
						}
						if p21 != 0 {
							savearchivo("camion2.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion2#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p21 = 0
						}
					}
				}
				//fmt.Println("Camion 2, Paquete: entregado")
				////fmt.Println("Paquete en la cola")
			}
			capacidad2 = 2
			intentos = 1
			intentos2 = 1

		}
	}()
	go func() { //camion 3
		capacidad3 := 2
		random := 0
		intentos := 1
		intentos2 := 1
		var id1 string
		var tipo1 string
		var valor1 string
		var origen1 string
		var destino1 string
		var id2 string
		var tipo2 string
		var valor2 string
		var origen2 string
		var destino2 string
		//log.Printf("Camion 3 listo")
		for {
			//time.Sleep(200 * time.Millisecond)
			if capacidad3 == 2 {
				chdisp3 <- "ready"
			}
			//time.Sleep(time.Second * 2)

			select {
			case <-ch3:
				//fmt.Println("Camion 3 en la cola")
				//se recibe un paquete
				if <-chdisp3 == "yes" {
					capacidad3--
					//fmt.Println("Camion 3 espera 2do paquete")
				}
				//pedir
				select {
				case <-chp3:
					respuesta1, err := p.PedirPaquete(context.Background(), &chat.Message{Mensaje: tipo31})
					if err != nil {
						log.Fatalf("Error al pedir paquete: %s", err)
					}
					log.Printf("Respuesta de logistica: %s", respuesta1.ID)
					p31, err = strconv.Atoi(respuesta1.Valor)
					id1 = respuesta1.ID
					tipo1 = respuesta1.Tipo
					valor1 = respuesta1.Valor
					origen1 = respuesta1.Origen
					destino1 = respuesta1.Destino
					_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(0) + "#En camino"})
					if err != nil {
						log.Fatalf("Error al actualizar paquete: %s", err)
					}
				}
				//time.Sleep(200 * time.Millisecond)
				select {
				case <-chdisp3:
					capacidad3--
					//pedir2
					select {
					case <-chp3:
						//fmt.Println("pidiendo 2do paquete recibido;", tipo32)
						if retail >= 1 || prioritario >= 1 || normal >= 1 {
							respuesta2, err := p.PedirPaquete(context.Background(), &chat.Message{Mensaje: tipo32})
							if err != nil {
								log.Fatalf("Error al pedir paquete: %s", err)
							}
							log.Printf("Respuesta de logistica: %s", respuesta2.ID)
							p32, err = strconv.Atoi(respuesta2.Valor)
							id2 = respuesta2.ID
							tipo2 = respuesta2.Tipo
							valor2 = respuesta2.Valor
							origen2 = respuesta2.Origen
							destino2 = respuesta2.Destino
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion2#" + strconv.Itoa(0) + "#En camino"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							//fmt.Println("2do paquete recibido")
						}
					}
				case <-time.After(time.Duration(espera) * time.Second): //tiempo de espera por 2do paquete
					//fmt.Println("no llegó 2do pedido")
				}
				//fmt.Println("Camion 3 se va")

				//funcion para entregar paquetes
				//se debe ver cual genera mas ingreso
				switch {
				case p32 == 0: //lego solo 1 paquete
					switch {
					case tipo31 == "Normal":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p31-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p31-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
						}
						if p31 != 0 {
							savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p31 = 0
						}
					case tipo31 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p31-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p31-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
						}
						if p31 != 0 {
							savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p31 = 0
						}
					}
				case p31 > p32 && p32 != 0: //el 1er paquete vale más
					switch {
					case tipo31 == "Normal":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p31-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p31-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
						}
						if p31 != 0 {
							savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p31 = 0
						}
					case tipo31 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p31-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p31-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
						}
						if p31 != 0 {
							savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p31 = 0
						}
					}
					switch {
					case tipo32 == "Normal":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3
								savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p32 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p32-10 > 0:
								intentos2++
								//fmt.Println("Reintento")
							case p32-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p32 = 0
								break
							}
						}
						if p32 != 0 {
							savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p32 = 0
						}
					case tipo32 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3(paquete2,intentos,fecha)
								savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p32 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p32-10 > 0:
								intentos2++
								//fmt.Println("Reintento")
							case p32-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p32 = 0
								break
							}
						}
						if p32 != 0 {
							//ESCRIBIR ARCHIVO CAMION3(paquete2,intentos,fecha)
							savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p32 = 0
						}
					}
				case p32 > p31: //el 2do paquete vale mas
					switch {
					case tipo32 == "Normal":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3(paquete2,intentos,fecha)
								savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p32 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p32-10 > 0:
								intentos2++
								//fmt.Println("Reintento")
							case p32-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p32 = 0
								break
							}
						}
						if p32 != 0 {
							//ESCRIBIR ARCHIVO CAMION3(paquete2,intentos,fecha)
							savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p32 = 0
						}
					case tipo32 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3(paquete2,intentos,fecha)
								savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p32 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p32-10 > 0:
								intentos2++
								//fmt.Println("Reintento")
							case p32-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p32 = 0
								break
							}
						}
						if p32 != 0 {
							//ESCRIBIR ARCHIVO CAMION3(paquete2,intentos,fecha)
							savearchivo("camion3.csv", id2, tipo2, valor2, origen2, destino2, intentos2, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id2 + "#camion3#" + strconv.Itoa(intentos2) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p32 = 0
						}
					}
					switch {
					case tipo31 == "Normal":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3(paquete1,intentos,fecha)
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p31-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p31-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
						}
						if p31 != 0 {
							//ESCRIBIR ARCHIVO CAMION3(paquete1,intentos,fecha)
							savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p31 = 0
						}
					case tipo31 == "Prioritario":
						for reintento := 0; reintento < 3; reintento++ {
							time.Sleep(time.Second * time.Duration(entrega))
							random = rand.Intn(5)
							if random < 4 {
								//exitoso
								//ESCRIBIR ARCHIVO CAMION3(paquete1,intentos,fecha)
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 0)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
							//fmt.Println("intento fallido")
							switch {
							case p31-10 > 0:
								intentos++
								//fmt.Println("Reintento")
							case p31-10 < 0:
								//fmt.Println("No lo vale")
								savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
								_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
								if err != nil {
									log.Fatalf("Error al actualizar paquete: %s", err)
								}
								p31 = 0
								break
							}
						}
						if p31 != 0 {
							//ESCRIBIR ARCHIVO CAMION3(paquete1,intentos,fecha)
							savearchivo("camion3.csv", id1, tipo1, valor1, origen1, destino1, intentos, 1)
							_, err = a.ActualizarPaquete(context.Background(), &chat.Message{Mensaje: id1 + "#camion3#" + strconv.Itoa(intentos) + "#No Entregado"})
							if err != nil {
								log.Fatalf("Error al actualizar paquete: %s", err)
							}
							p31 = 0
						}
					}
				}
				//fmt.Println("Camion 3, Paquete entregado")

				//fmt.Println("Paquete en la cola")
			}
			capacidad3 = 2
			intentos = 1
			intentos2 = 1
		}

	}()
	capacidad := 2 //representa la capacidad actual del camion que se esta cargando
	//retail := 3
	//prioritario := 3
	//normal := 3
	for {
		response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
		if err != nil {
			log.Fatalf("aaaah Error al verificar si hay pedidos: %s", err)
		}
		log.Printf("Respuesta de logistica: %s", response.Respuesta)
		//time.Sleep(time.Second * 2)
		// Tengo el primer paquete
		if response.Respuesta == "Si" {
			retail, err := strconv.Atoi(response.Retail)
			if err != nil {
				log.Fatalf("Error al verificar si hay pedidos(retail):", err)
			}
			prioritario, err := strconv.Atoi(response.Prioritario)
			if err != nil {
				log.Fatalf("Error al verificar si hay pedidos(prioritario):", err)
			}
			normal, err := strconv.Atoi(response.Normal)
			if err != nil {
				log.Fatalf("Error al verificar si hay pedidos(normal):", err)
			}
			fmt.Println("Cantidad retail:", retail)
			fmt.Println("Cantidad prioritario:", prioritario)
			fmt.Println("Cantidad normal:", normal)
			if retail >= 1 || prioritario >= 1 { //hay paquetes en cola
				//se revisa camion 1
				select {
				case <-chdisp1:
					capacidad = 2 //capacidad actual del camion
					//si esta listo se le envian 2 paquete
					ch <- "go1"
					//time.Sleep(100 * time.Millisecond)
					if retail >= 1 {
						chdisp1 <- "yes"
						//fmt.Println("se envio un retail a camion 1")
						capacidad--
						tipo11 = "Retail"
						chp <- "ok"
						if retail > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("1.Error al verificar si hay pedidos: %s", err)
							}
							if response.Retail == "" {
								retail = 0
							} else {
								retail, err = strconv.Atoi(response.Retail)
								if err != nil {
									log.Fatalf("Error al verificar si hay pedidos(retail):", err)
								}
							}
						} else {
							retail = 0
						}
					}
					//time.Sleep(100 * time.Millisecond)
					//se envia 1 si aun quedan paquetes
					if retail >= 1 {
						chdisp1 <- "yes"
						//fmt.Println("se envio un retail a camion 1")
						capacidad--
						switch {
						case p11 == 0:
							tipo11 = "Retail"
						case p11 > 0:
							tipo12 = "Retail"
						}
						chp <- "ok"
						if retail > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("2.Error al verificar si hay pedidos: %s", err)
							}
							if response.Retail == "" {
								retail = 0
							} else {
								retail, err = strconv.Atoi(response.Retail)
								if err != nil {
									log.Fatalf("Error al verificar si hay pedidos(retail):", err)
								}
							}
						} else {
							retail = 0
						}
					}
					//time.Sleep(100 * time.Millisecond)
					if prioritario >= 1 && capacidad >= 1 {
						chdisp1 <- "yes"
						//fmt.Println("se envio prioritario a camion 1")
						capacidad--
						switch {
						case p11 == 0:
							tipo11 = "Prioritario"
						case p11 > 0:
							tipo12 = "Prioritario"
						}
						chp <- "ok"
						if prioritario > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("3.Error al verificar si hay pedidos: %s", err)
							}
							prioritario, err = strconv.Atoi(response.Prioritario)
							if err != nil {
								log.Fatalf("Error al verificar si hay pedidos(prioritario):", err)
							}
						} else {
							prioritario = 0
						}
					}
					//time.sleep(100 * time.Millisecond)
					if prioritario >= 1 && capacidad >= 1 {
						chdisp1 <- "yes"
						//fmt.Println("se envio prioritario a camion 1")
						capacidad--
						switch {
						case p11 == 0:
							tipo11 = "Prioritario"
						case p11 > 0:
							tipo12 = "Prioritario"
						}
						chp <- "ok"
						if prioritario > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("4.Error al verificar si hay pedidos: %s", err)
							}
							prioritario, err = strconv.Atoi(response.Prioritario)
							if err != nil {
								log.Fatalf("Error al verificar si hay pedidos(prioritario):", err)
							}
						} else {
							prioritario = 0
						}
					}
				case <-time.After(1 * time.Second):
					//fmt.Println("Camion 1 ocupado")
				}
			}
			if retail >= 1 || prioritario >= 1 { //hay paquetes en cola
				//se revisa camion 2
				select {
				case <-chdisp2:
					capacidad = 2 //capacidad actual del camion
					//si esta listo se le envian 2 paquete
					ch2 <- "go2"
					//time.sleep(100 * time.Millisecond)
					if retail >= 1 {
						chdisp2 <- "yes"
						//fmt.Println("se envio un retail a camion 2")
						capacidad--
						tipo21 = "Retail"
						chp2 <- "ok"
						if retail > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("5.Error al verificar si hay pedidos: %s", err)
							}
							if response.Retail == "" {
								retail = 0
							} else {
								retail, err = strconv.Atoi(response.Retail)
								if err != nil {
									log.Fatalf("Error al verificar si hay pedidos(retail):", err)
								}
							}
						} else {
							retail = 0
						}
					}
					//time.sleep(100 * time.Millisecond)
					//se envia 1 si aun quedan paquetes
					if retail >= 1 {
						chdisp2 <- "yes"
						//fmt.Println("se envio un retail a camion 2")
						capacidad--
						switch {
						case p21 == 0:
							tipo21 = "Retail"
						case p21 > 0:
							tipo22 = "Retail"
						}
						chp2 <- "ok"
						if retail > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("6.Error al verificar si hay pedidos: %s", err)
							}
							if response.Retail == "" {
								retail = 0
							} else {
								retail, err = strconv.Atoi(response.Retail)
								if err != nil {
									log.Fatalf("Error al verificar si hay pedidos(retail):", err)
								}
							}
						} else {
							retail = 0
						}
					}
					//time.sleep(100 * time.Millisecond)
					if prioritario >= 1 && capacidad >= 1 {
						chdisp2 <- "yes"
						//fmt.Println("se envio prioritario a camion 2")
						capacidad--
						switch {
						case p21 == 0:
							tipo21 = "Prioritario"
						case p21 > 0:
							tipo22 = "Prioritario"
						}
						chp2 <- "ok"
						if prioritario > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("7.Error al verificar si hay pedidos: %s", err)
							}
							prioritario, err = strconv.Atoi(response.Prioritario)
							if err != nil {
								log.Fatalf("Error al verificar si hay pedidos(prioritario):", err)
							}
						} else {
							prioritario = 0
						}
					}
					//time.sleep(100 * time.Millisecond)
					if prioritario >= 1 && capacidad >= 1 {
						chdisp2 <- "yes"
						//fmt.Println("se envio prioritario a camion 2")
						capacidad--
						switch {
						case p21 == 0:
							tipo21 = "Prioritario"
						case p21 > 0:
							tipo22 = "Prioritario"
						}
						chp2 <- "ok"
						if prioritario > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("8.Error al verificar si hay pedidos: %s", err)
							}
							prioritario, err = strconv.Atoi(response.Prioritario)
							if err != nil {
								log.Fatalf("Error al verificar si hay pedidos(prioritario):", err)
							}
						} else {
							prioritario = 0
						}
					}
				case <-time.After(1 * time.Second):
					//fmt.Println("Camion 2 ocupado")
				}
			}
			//time.sleep(100 * time.Millisecond)
			if normal >= 1 || prioritario >= 1 {
				//se revisa camion 3
				select {
				case <-chdisp3:
					capacidad = 2 //capacidad actual del camion
					//si esta listo se le envian 2 paquete
					ch3 <- "go3"
					//time.sleep(100 * time.Millisecond)
					if prioritario >= 1 {
						chdisp3 <- "yes"
						//fmt.Println("se envio prioritario a camion 3")
						capacidad--
						tipo31 = "Prioritario"
						chp3 <- "ok"
						if prioritario > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("9.Error al verificar si hay pedidos: %s", err)
							}
							prioritario, err = strconv.Atoi(response.Prioritario)
							if err != nil {
								log.Fatalf("Error al verificar si hay pedidos(prioritario):", err)
							}
						} else {
							prioritario = 0
						}
					}
					//time.sleep(100 * time.Millisecond)
					if prioritario >= 1 {
						chdisp3 <- "yes"
						//fmt.Println("se envio prioritario a camion 3")
						capacidad--
						switch {
						case p31 == 0:
							tipo31 = "Prioritario"
						case p31 > 0:
							tipo32 = "Prioritario"
						}
						chp3 <- "ok"
						if prioritario > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("10.Error al verificar si hay pedidos: %s", err)
							}
							prioritario, err = strconv.Atoi(response.Prioritario)
							if err != nil {
								log.Fatalf("Error al verificar si hay pedidos(prioritario):", err)
							}
						} else {
							prioritario = 0
						}
					}
					//time.sleep(100 * time.Millisecond)
					if normal >= 1 && capacidad >= 1 {
						chdisp3 <- "yes"
						//fmt.Println("se envio un normal a camion 3")
						capacidad--
						switch {
						case p31 == 0:
							tipo31 = "Normal"
						case p31 > 0:
							tipo32 = "Normal"
						}
						chp3 <- "ok"
						if normal > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("11.Error al verificar si hay pedidos: %s", err)
							}
							normal, err = strconv.Atoi(response.Normal)
							if err != nil {
								log.Fatalf("Error al verificar si hay pedidos(normal):", err)
							}
						} else {
							normal = 0
						}
					}
					//time.sleep(100 * time.Millisecond)
					if normal >= 1 && capacidad >= 1 {
						chdisp3 <- "yes"
						//fmt.Println("se envio un normal a camion 3")
						capacidad--
						switch {
						case p31 == 0:
							tipo31 = "Normal"
						case p31 > 0:
							tipo32 = "Normal"
						}
						chp3 <- "ok"
						if normal > 1 {
							response, err := c.VerificarPedido(context.Background(), &chat.Message{Mensaje: "Tienes pedidos en las colas"})
							if err != nil {
								log.Fatalf("12.Error al verificar si hay pedidos: %s", err)
							}
							normal, err = strconv.Atoi(response.Normal)
							if err != nil {
								log.Fatalf("Error al verificar si hay pedidos(normal):", err)
							}
						} else {
							normal = 0
						}
					}
				case <-time.After(1 * time.Second):
					//fmt.Println("Camion 3 ocupado")
				}
			}

			time.Sleep(100 * time.Millisecond)

			//retail, err := strconv.Atoi(response.Retail)
			//if err != nil {
			//	log.Fatalf("Error al pasar el string retail a entero: %s", err)
			//}

			log.Printf("Tenemos pedidos en las colas")

		}
		//por el momento termina el for si tenemos paquetes en las colas
		//time.Sleep(time.Second * 3)
		//break
	}
}
