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
	"time"

	"google.golang.org/grpc"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
)

// Nodes -> DataNodes
var Nodes = []string{":9001", ":9002", ":9003"}

// random -> [0,1,2]
func random() int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	random := r.Intn(3)
	return random
}

func split() {

	random := random()
	//conn, err := grpc.Dial(Nodes[random], grpc.WithInsecure())

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(Nodes[random], grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	l := chat.NewLibroServiceClient(conn)
	c := chat.NewChunkServiceClient(conn)
	p := chat.NewPropuestaServiceClient(conn)

	fmt.Print("Nombre de libro a separar: ")
	var libro string
	fmt.Scanln(&libro)
	//fmt.Println()

	fileToBeChunked := "./" + libro + ".pdf"

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

	// pasar el random a string
	_, err = p.GenProp(context.Background(), &chat.Node{IDNode: strconv.Itoa(random)})
	if err != nil {
		log.Fatalf("Error al enviar chunk")
	}
}

func recombine(libro string, totalPartsNum uint64) {

	//fmt.Println()

	fileToBeChunked := "./" + libro + ".pdf"
	//fileToBeChunked := "./bigfile.zip" // change here!

	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()
	//const fileChunk = 1 * (1 << 20) // 1 MB, change this to your requirement

	// calculate total number of parts the file will be chunked into

	//totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	// just for fun, let's recombine back the chunked files in a new file

	newFileName := "NEW" + libro + ".pdf"
	_, err = os.Create(newFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//set the newFileName file to APPEND MODE!!
	// open files r and w

	file, err = os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

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
		currentChunkFileName := libro + "_" + strconv.FormatUint(j, 10)

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

	fmt.Println("\nComportamiento a seguir (Uploader/Downloader): ")
	var comportamiento string
	fmt.Scanln(&comportamiento)

	if comportamiento == "Uploader" {

		// Crear un random para ver a que nodo enviar las cosas
		// Enviar nombre del libro, cantidad de partes y chunks
		// Enviar un mensaje diciendo que ya se envio todo

		//Generar numero random
		split()

	}
}
