package product

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type EventData struct {
	ID          *string  `json:"product_id"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	AccountID   *int     `json:"accountID"`
}

type Event struct {
	Type string    `json:"type"`
	Data EventData `json:"data"`
}

var done = make(chan bool)

func (service productService) SendMessageToRecommender(event Event, topic string) error {
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

func (service productService) MsgHandler() {
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

func (service productService) Close() {
	if err := service.producer.Close(); err != nil {
		log.Printf("Failed to close producer: %v\n", err)
	} else {
		done <- true
	}
}
