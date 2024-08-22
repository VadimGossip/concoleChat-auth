package producer

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

type producer struct {
	syncProducer sarama.SyncProducer
}

func NewProducer(syncProducer sarama.SyncProducer) *producer {
	return &producer{
		syncProducer: syncProducer,
	}
}

func (p *producer) Produce(topic string, msg any) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %s", err)
	}

	_, _, err = p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *producer) Close() error {
	return p.syncProducer.Close()
}
