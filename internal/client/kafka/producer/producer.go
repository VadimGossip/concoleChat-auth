package producer

import (
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

type marshaller interface {
	Marshal() ([]byte, error)
}

func (p *producer) Produce(topic string, msg any) error {
	if val, ok := msg.(marshaller); ok {
		data, err := val.Marshal()
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

	return fmt.Errorf("the message must implement the method Marshal() ([]byte, error)")
}

func (p *producer) Close() error {
	return p.syncProducer.Close()
}
