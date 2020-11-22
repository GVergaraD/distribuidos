package chat

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Server -> Estructura
type Server struct {
}

var next *list.Element

// Nodes -> DataNodes (1,2,3)
var Nodes = []string{":9001", ":9002", ":9003"}

// PartesNode -> variable para guardar la cantidad de partes del libro
var PartesNode [3]int

// TituloNode -> variable para guardar el titulo del libro
var TituloNode [3]string

// ChunksNode1 -> variable para guardar los chunks del libro para Node1
var ChunksNode1 = list.New()

// ChunksNode2 -> variable para guardar los chunks del libro para Node2
var ChunksNode2 = list.New()

// ChunksNode3 -> variable para guardar los chunks del libro para Node3
var ChunksNode3 = list.New()

// ChunksNode -> Lista de todos los chunks para generalizar
var ChunksNode = []*list.List{ChunksNode1, ChunksNode2, ChunksNode3}

// ElegirNode -> recibe un conjunto de nodos y escoge uno de ellos aleatoriamente
func ElegirNode(participantes []int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	random := r.Intn(len(participantes))
	return participantes[random]
}

// distribuir -> funcion que recibe un chunk y lo distribuye a "node"
func distribuir(fileName string, node string) {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(node, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewChunkServiceClient(conn)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error al leer el archivo")
	}

	_, err = c.PasarChunk(context.Background(), &Chunks{ID: "3", FileName: fileName, Chunk: data})
	if err != nil {
		log.Fatalf("Error al distribuir el chunk %s", fileName)
	}
}

// contactar -> funcion que revisa si el nodo acepta la propuesta
func contactar(node string) bool {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(node, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewContacarServiceClient(conn)

	_, err = c.ContacarNode(context.Background(), &Message{Mensaje: "Aceptas la propuesta?"})
	if err != nil {
		return true // Nodo no responde
	}

	return false // Nodo si responde
}

// Ints returns a unique subset of the int slice provided.
func Ints(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

// PasarLibro -> funcion que recibe el titulo y la cantidad de partes del libro desde el cliente
func (s *Server) PasarLibro(ctx context.Context, in *Libro) (*MessageResponse, error) {

	log.Printf("Libro del cliente: %s", in.Nombre)
	log.Printf("Cantidad de partes del libro del cliente: %s", in.Partes)

	num, err := strconv.Atoi(in.Partes)
	if err != nil {
		log.Fatalf("Error al pasar string a entero")
	}

	ID, err := strconv.Atoi(in.ID)
	if err != nil {
		log.Fatalf("Error al pasar string a entero")
	}

	TituloNode[ID] = in.Nombre
	PartesNode[ID] = num

	return &MessageResponse{Respuesta: "Libro recibido"}, nil
}

// PasarChunk -> funcion que recibe los chunks desde el cliente
func (s *Server) PasarChunk(ctx context.Context, in *Chunks) (*MessageResponse, error) {

	log.Printf("Chunk del cliente: %s", in.FileName)

	_, err := os.Create(in.FileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// write/save buffer to disk
	ioutil.WriteFile(in.FileName, in.Chunk, os.ModeAppend)

	IDNod, err := strconv.Atoi(in.ID) // in.ID = "0" or "1" or "2" or "3" | 3 -> No hay que guardar el chunk, solo se esta distribuyendo
	if err != nil {
		log.Fatalf("Error al pasar string a entero")
	}

	if IDNod != 3 {
		ChunksNode[IDNod].PushBack(in.FileName)
	}

	return &MessageResponse{Respuesta: "Chunk recibido"}, nil
}

// ContacarNode -> funcion para ver si el nodo acepta la propuesta
func (s *Server) ContacarNode(ctx context.Context, in *Message) (*MessageResponse, error) {

	return &MessageResponse{Respuesta: "Propuesta aceptada"}, nil
}

// GenProp -> funcion que genera la propuesta y verfica que esta sea aceptada,
// si no es aceptada genera una nueva hasta converger
func (s *Server) GenProp(ctx context.Context, in *Node) (*MessageResponse, error) {

	IDNodo, err := strconv.Atoi(in.IDNode) // in.IDNode = "0" or "1" or "2"
	if err != nil {
		log.Fatalf("Error al pasar string a entero")
	}

	// Crear lista y guardar los chunks a distribuir
	for e := ChunksNode[IDNodo].Front(); e != nil; e = next {
		next = e.Next()
		log.Printf("%s", e.Value)
	}

	band := true

	//Nodos que van a participar en la propuesta inicial
	var nodos = []int{0, 1, 2}

	// Propuesta a generar
	propuesta := make([]int, ChunksNode[IDNodo].Len())

	for band {

		if PartesNode[IDNodo] >= 3 {

			// En principio reparte 1 chunk a cada nodo y el resto lo reparte de forma aleatoria
			for i := 0; i < len(nodos); i++ {
				propuesta[i] = nodos[i]
				// Al principio tenemos 3 nodos
				// -> propuesta[0] = primer nodo de la lista
				// -> propuesta[1] = segundo nodo de la lista
				// -> propuesta[2] = tercer nodo de la lista
			}
			//enviar de forma random las siguientes
			for i := len(nodos); i < ChunksNode[IDNodo].Len(); i++ {
				propuesta[i] = ElegirNode(nodos)
				// -> propuesta[3] = random
				// ...
			}
		} else {
			// Se queda el con una y la otra la reparte de forma aleatoria
			propuesta[0] = IDNodo

			//enviar de forma random la siguiente
			var NodosAux = nodos
			NodosAux = remove(NodosAux, IDNodo)
			propuesta[1] = ElegirNode(NodosAux)
		}

		// Participantes de la propuesta
		unique := Ints(propuesta)

		// Preguntar si la propuesta ha sido aceptada por cada nodo
		for i := 0; i < len(unique); i++ {
			if unique[i] != IDNodo {
				// Si todo responden band = false
				if contactar(Nodes[unique[i]]) == true {
					// unique[i] -> Nodo no respondio, hay que sacarlo de la lista y generar otra propuesta
					nodos = remove(nodos, unique[i])
					band = true
					break
				} else { // Si todo responden band = false -> puedo mandar la propuesta y distribuir
					band = false
				}
			}
		}
	}
	// Escribir en el Log (Crear a NameNode)

	// Distribuir propuesta
	j := 0
	for e := ChunksNode[IDNodo].Front(); e != nil; e = next {
		next = e.Next()
		if propuesta[j] != IDNodo {
			distribuir(e.Value.(string), Nodes[propuesta[j]])
			// Eliminar el archivo
			// Activar la siguiente linea solo en ubuntu
			os.Remove(e.Value.(string))
		}
		ChunksNode[IDNodo].Remove(e)
		j++
	}

	return &MessageResponse{Respuesta: "Propuesta generada"}, nil
}
