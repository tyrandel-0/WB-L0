package main

import (
	"github.com/nats-io/stan.go"
	"io"
	"log"
	"os"
)

func main() {
	path := "data"
	channel := "test"

	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}

	defer func(sc stan.Conn) {
		err := sc.Close()
		if err != nil {
			return
		}
	}(sc)

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		msg, err := readData(path + "/" + file.Name())
		if err != nil {
			return
		}
		if err = sc.Publish(channel, msg); err != nil {
			log.Fatalf("Error publishing: %v", err)
		} else {
			log.Printf("Published from: %s", file.Name())
		}
	}
}

func readData(path string) ([]byte, error) {
	file, err := os.Open(path)
	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
