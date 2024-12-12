package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Order struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

func main() {
	// Kết nối NATS
	nc, err := nats.Connect("nats://192.168.1.249:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Tạo JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Tạo stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"orders.*"},
		Storage:  nats.FileStorage,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Publisher goroutine
	go func() {
		for i := 0; i < 10; i++ {
			order := Order{
				ID:     fmt.Sprintf("order-%d", i),
				Amount: float64(i * 100),
			}
			data, _ := json.Marshal(order)
			_, err := js.Publish("orders.new", data)
			if err != nil {
				log.Printf("Error publishing: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()

	// Consumer
	_, err = js.Subscribe("orders.*", func(msg *nats.Msg) {
		var order Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Printf("Error unmarshaling: %v", err)
			msg.Nak()
			return
		}

		log.Printf("Received order: %+v", order)
		msg.Ack()
	}, nats.Durable("orders-processor"))
	js.Subscribe(">", func(msg *nats.Msg) {
		fmt.Printf("Subject: %s\nData: %s\n", msg.Subject, string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}

	// Wait forever
	select {}
}
