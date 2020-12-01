package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
)

// Nodes -> DataNodes
var Nodes = []string{"10.6.40.221:9001", "10.6.40.222:9002", "10.6.40.223:9003"}

// random -> [0,1,2]
func random() int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	random := r.Intn(3)
	return random
}

func pedir(filename string, node string) {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(node, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	p := chat.NewChunkServiceClient(conn)

	_, err = p.PedirChunk(context.Background(), &chat.Chunks{FileName: filename})
	if err != nil {
		log.Printf("Error al consutar a NameNode")
	}

}

func split(mode string) {

	random := random()
	random = 0

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(Nodes[random], grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	l := chat.NewLibroServiceClient(conn)
	c := chat.NewChunkServiceClient(conn)
	p := chat.NewPropuestaServiceClient(conn)

	fmt.Print("Nombre de libro a subir: ")
	var libro string
	fmt.Scanln(&libro)
	//fmt.Println()

	fileToBeChunked := "./Libros/" + libro + ".pdf"

	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	//const fileChunk = 1 * (1 << 20) // 1 MB, change this to your requirement
	const fileChunk = 256000
	fmt.Printf("Valor de filechunk: %d .\n", fileChunk)

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	TotalParts := strconv.FormatUint(totalPartsNum, 10)

	response, err := l.PasarLibro(context.Background(), &chat.Libro{ID: strconv.Itoa(random), Nombre: libro, Partes: TotalParts})
	if err != nil {
		log.Fatalf("Error al enviar libro")
	}
	log.Printf("Libro ingresado, respuesta del nodo: %s", response.Respuesta)

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// write to disk
		fileName := libro + "_" + strconv.FormatUint(i, 10)

		response, err := c.PasarChunk(context.Background(), &chat.Chunks{ID: strconv.Itoa(random), FileName: fileName, Chunk: partBuffer})
		if err != nil {
			log.Fatalf("Error al enviar chunk")
		}
		log.Printf("Chunk ingresado, respuesta del nodo: %s", response.Respuesta)

		fmt.Println("Split to : ", fileName)
	}

	// Avisa que se terminaron de pasar los chunks en el DataNode para que este genere la propuesta

	if mode == "Distribuido" { // Generar la propuesta "Exclusion Mutua Distribuida"
		_, err = p.GenProp(context.Background(), &chat.Node{IDNode: strconv.Itoa(random)})
		if err != nil {
			log.Fatalf("Error al enviar chunk")
		}
	}
	if mode == "Centralizado" { // Generar la propuesta "Exclusion Mutua Centralizada"
		_, err = p.GenProp2(context.Background(), &chat.Node{IDNode: strconv.Itoa(random)})
		if err != nil {
			log.Fatalf("Error al enviar chunk")
		}
	}
}

func recombine(libro string, totalPartsNum uint64) {

	newFileName := "Descargado_de_Vuamos_" + libro + ".pdf"
	_, err := os.Create(newFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//set the newFileName file to APPEND MODE!!
	// open files r and w

	file, err := os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
	// defer file.Close()

	// just information on which part of the new file we are appending
	var writePosition int64 = 0

	for j := uint64(0); j < totalPartsNum; j++ {

		//read a chunk
		currentChunkFileName := "../chunks/" + libro + "_" + strconv.FormatUint(j, 10)

		newFileChunk, err := os.Open(currentChunkFileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer newFileChunk.Close()

		chunkInfo, err := newFileChunk.Stat()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// calculate the bytes size of each chunk
		// we are not going to rely on previous data and constant

		var chunkSize int64 = chunkInfo.Size()
		chunkBufferBytes := make([]byte, chunkSize)

		fmt.Println("Appending at position : [", writePosition, "] bytes")
		writePosition = writePosition + chunkSize

		// read into chunkBufferBytes
		reader := bufio.NewReader(newFileChunk)
		_, err = reader.Read(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
		// write/save buffer to disk
		//ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)

		n, err := file.Write(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file.Sync() //flush to disk

		// free up the buffer for next cycle
		// should not be a problem if the chunk size is small, but
		// can be resource hogging if the chunk size is huge.
		// also a good practice to clean up your own plate after eating

		chunkBufferBytes = nil // reset or empty our buffer

		fmt.Println("Written ", n, " bytes")

		fmt.Println("Recombining part [", j, "] into : ", newFileName)
	}

	// now, we close the newFileName
	file.Close()

}

func main() {

	fmt.Println("Modo del programa(Distribuido/Centralizado): ")
	var modo string
	fmt.Scanln(&modo)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.6.40.224:9004", grpc.WithInsecure()) // Para conectarse a NameNode
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewConsultarServiceClient(conn)

	var value int

	for {
		fmt.Println("\nIgrese opcion a realizar")
		fmt.Println("1.- Upload")
		fmt.Println("2.- Download")
		fmt.Printf("3.- Salir\n\n")
		fmt.Scanln(&value)

		if value == 1 {

			start := time.Now()
			split(modo)
			elapsed := time.Since(start)
			fmt.Println("Tiempo en subir el archivo:", elapsed)

		} else if value == 2 {

			for {
				fmt.Println("\nQue opcion desea realizar: ")
				fmt.Println("1) Ver listado de libros disponibles")
				fmt.Println("2) Descargar libro")
				fmt.Printf("3) Volver\n\n")
				var option int
				fmt.Scanln(&option)

				if option == 1 { //Pedir los libros disponibles

					respuesta, err := c.Consultar(context.Background(), &chat.Consult{Tipo: "Libros"})
					if err != nil {
						log.Printf("Error al consutar a NameNode")
					}

					fmt.Println("\nLibros disponibles:\n")
					for i := 0; i < len(respuesta.Respuesta); i++ {
						fmt.Println(strconv.Itoa(i+1)+")", respuesta.Respuesta[i])
					}
					fmt.Println("")

				} else if option == 2 { //Descargar libro

					fmt.Println("\nIngrese libro a descargar: ")
					var libro string
					fmt.Scanln(&libro)

					// Pedimos la info de los chunks al NameNode
					respuesta, err := c.Consultar(context.Background(), &chat.Consult{Tipo: "Ubicacion", Titulo: libro})
					if err != nil {
						log.Printf("Error al consutar a NameNode")
					}

					// Variables para guardar los chunks y su ubicacion
					chunks := make([]string, 0, 100)
					ubicaciones := make([]string, 0, 100)

					// Guardamos los chunks y su ubicacion
					for i := 0; i < len(respuesta.Respuesta); i++ {
						respuesta := strings.Split(respuesta.Respuesta[i], " ")
						chunk, ubicacion := respuesta[0], respuesta[1]
						chunks = append(chunks, chunk)
						ubicaciones = append(ubicaciones, ubicacion)
					}

					//Pedir a los DataNode los chunks en base a su ubicacion
					for j := 0; j < len(chunks); j++ {
						fmt.Println(chunks[j], ubicaciones[j])
						pedir(chunks[j], ubicaciones[j])
					}

					//Reconstruir el libro
					recombine(libro, uint64(len(chunks)))

				} else if option == 3 {
					break
				}
			}

		} else if value == 3 {
			break
		}
	}
}
