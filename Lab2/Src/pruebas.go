package main

import (
	"fmt"
	"math/rand"
	"time"
)

func random2(participantes []int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	random := r.Intn(len(participantes))
	return participantes[random]
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func main() {

	var nodos = []int{0, 1, 2}
	//random := random2(nodos)
	fmt.Println(nodos)
	nodos = remove(nodos, 0)
	fmt.Println(nodos)
}
