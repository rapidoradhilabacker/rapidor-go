package pulsar

import (
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

const (
	pulsarURL = "pulsar://localhost:6650" // Adjust if needed
	topicName = "persistent://public/default/my-topic"
)

type PulsarClient struct {
	client   pulsar.Client
	producer pulsar.Producer
}

func NewPulsarClient() *PulsarClient {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: pulsarURL,
	})
	if err != nil {
		log.Fatalf("Could not create Pulsar client: %v", err)
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topicName,
	})
	if err != nil {
		log.Fatalf("Could not create Pulsar producer: %v", err)
	}

	return &PulsarClient{client: client, producer: producer}
}

func (p *PulsarClient) Produce(message []byte) error {
	_, err := p.producer.Send(pulsar.ProducerMessage{
		Payload: message,
	})
	if err != nil {
		log.Printf("Failed to publish message to Pulsar: %v", err)
	}
	return err
}

func (p *PulsarClient) Close() {
	p.producer.Close()
	p.client.Close()
}
