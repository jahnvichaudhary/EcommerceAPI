package product

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
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

	// Optionally, handle errors and successes
	go func() {
		for {
			select {
			case success := <-service.producer.Successes():
				log.Printf("Message sent to partition %d at offset %d\n", success.Partition, success.Offset)
			case err = <-service.producer.Errors():
				log.Printf("Failed to send message: %v\n", err)
			}
		}
	}()

	return nil
}
