package main

import (
	"github.com/nats-io/stan.go"
	"log"
	"os"
)

func main() {
	sc, err := stan.Connect("test-cluster", "pub")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	if len(os.Args) != 2 {
		log.Fatal("Put one filename")
	}
	filename := os.Args[1]

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	sc.Publish("order", file)

}
