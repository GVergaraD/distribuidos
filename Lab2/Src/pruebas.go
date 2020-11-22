package main

import (
	"container/list"
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

// ChunksNode1 wea
var ChunksNode1 = list.New()

// ChunksNode2 wea
var ChunksNode2 = list.New()

func main() {

	var nodos = []*list.List{ChunksNode1, ChunksNode2}
	//random := random2(nodos)
	ChunksNode1.PushBack(0)
	ChunksNode1.PushBack(1)
	ChunksNode1.PushBack(2)
	nodos[0].PushBack(3)
	fmt.Println(nodos[0].Len())
}
