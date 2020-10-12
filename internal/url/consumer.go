package url

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

type KafkaConsumer struct {
	c sarama.Consumer
}

func NewKafkaConsumer(host, port string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := []string{host + ":" + port}
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer}, nil
}

func (co *KafkaConsumer) Close() error {
	return co.c.Close()
}

func (co *KafkaConsumer) Consume(f func(*sarama.ConsumerMessage)) error {
	consumer, err := co.c.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err = <-consumer.Errors():
				log.Errorf("Consumer err: %v", err)
			case msg := <-consumer.Messages():
				f(msg)
			case <-signals:
				log.Info("System interrupt.")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh

	return nil
}
