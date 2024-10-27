package ws

// core packages are imported here
// ws is a package that contains the application's websocket logic
// ws is imported in other packages

import (
	"context"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// PulsarManager manages Pulsar connections and messaging
type PulsarManager struct {
	client pulsar.Client
}

// NewPulsarManager initializes a new Pulsar client
func NewPulsarManager(url string) (*PulsarManager, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               url,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
		return nil, err
	}
	return &PulsarManager{client: client}, nil
}

// PublishMessage publishes a message to the specified topic
func (pm *PulsarManager) PublishMessage(topic string, message []byte) error {
	producer, err := pm.client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return err
	}
	defer producer.Close()

	pulsarMessage := pulsar.ProducerMessage{
		Payload: message,
	}

	// Send the message
	_, err = producer.Send(context.Background(), &pulsarMessage)
	if err != nil {
		return err
	}

	log.Printf("Message sent to topic %s: %s", topic, string(message))
	return nil
}

func (pm *PulsarManager) Subscribe(topic string) (pulsar.Consumer, error) {
	consumer, err := pm.client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: "my-sub",
		Type:             pulsar.Shared,
	})
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// Close closes the Pulsar client
func (pm *PulsarManager) Close() {
	pm.client.Close()
}
