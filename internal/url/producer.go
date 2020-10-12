package url

import (
	"encoding/binary"
	"github.com/Shopify/sarama"
	"time"
)

const topic = "url_logs"

type Producer interface {
	Produce(int64, time.Time) error
}

type KafkaProducer struct {
	p sarama.SyncProducer
}

func NewKafkaProducer(host, port string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	brokers := []string{host + ":" + port}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer}, nil
}

func (pr *KafkaProducer) Produce(id int64, t time.Time) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(id))

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(b),
		Timestamp: t.UTC(),
	}

	_, _, err := pr.p.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (pr *KafkaProducer) Close() error {
	return pr.p.Close()
}
