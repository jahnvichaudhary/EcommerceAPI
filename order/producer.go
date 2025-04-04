package order

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type EventData struct {
	AccountId int    `json:"user_id"`
	ProductId string `json:"product_id"`
}

type Event struct {
	Type      string    `json:"type"`
	EventData EventData `json:"data"`
}

var done = make(chan bool)

func (service orderService) SendMessageToRecommender(event Event, topic string) error {
	jsonMessage, err := json.Marshal(event)
	if err != nil {
		log.Println("Failed to marshal event:", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMessage),
	}

	// Send the message asynchronously
	service.producer.Input() <- msg

	return nil
}

func (service orderService) MsgHandler() {
	go func() {
		for {
			select {
			case success := <-service.producer.Successes():
				log.Printf("Message sent to partition %d at offset %d\n", success.Partition, success.Offset)
			case err := <-service.producer.Errors():
				log.Printf("Failed to send message: %v\n", err)
			case <-done:
				log.Println("Producer closed successfully")
				return
			}
		}
	}()
}

func (service orderService) Close() {
	if err := service.producer.Close(); err != nil {
		log.Printf("Failed to close producer: %v\n", err)
	} else {
		done <- true
	}
}
