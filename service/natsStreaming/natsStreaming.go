package natsStreaming

import (
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"

	"service/server/cache"
	"service/server/http-server/model"
	"service/server/storage/postgres"
)

func RunNatsStreaming(storage *postgres.Storage, cache *cache.Cache) error {
	sc, err := stan.Connect("test-cluster", "subscriber-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer func() {
		if err := sc.Close(); err != nil {
			log.Printf(err.Error())
		}
	}()

	subject := "test"

	_, err = sc.Subscribe(subject, func(m *stan.Msg) {
		if err = handleMessage(storage, cache, m); err != nil {
			return
		}
	}, stan.StartWithLastReceived())
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	log.Println("Subscriber is listening...")

	select {}
}

func handleMessage(storage *postgres.Storage, cache *cache.Cache, m *stan.Msg) error {
	var order model.Model
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		return err
	}

	if err = storage.AddOrder(order); err != nil {
		return err
	}

	byt, err := json.Marshal(order)
	if err != nil {
		return err
	}
	cache.SetDefault(order.OrderUID, byt)

	log.Printf("Received a message: %s\n", order.OrderUID)

	return nil
}
