package producer

import "github.com/IBM/sarama"

type producer struct {
	syncProducer sarama.SyncProducer
}

func NewProducer(syncProducer sarama.SyncProducer) *producer {
	return &producer{
		syncProducer: syncProducer,
	}
}

func (p *producer) Produce(msg *sarama.ProducerMessage) error {
	_, _, err := p.syncProducer.SendMessage(msg)
	return err
}

func (p *producer) Close() error {
	return p.syncProducer.Close()
}
