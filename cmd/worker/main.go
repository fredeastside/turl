package main

import (
	"encoding/binary"
	"github.com/Shopify/sarama"
	"os"
	"turl/internal/db"
	"turl/internal/url"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)

	conn, err := db.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB conn error %s", err.Error())
	}
	defer conn.Close()

	c, err := url.NewKafkaConsumer(os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	if err != nil {
		log.Fatalf("Kafka conn error %s", err.Error())
	}
	defer c.Close()

	repo := url.NewURLRepository(conn, nil)
	err = c.Consume(func(msg *sarama.ConsumerMessage) {
		id := int64(binary.LittleEndian.Uint64(msg.Value))
		err = repo.Log(id, msg.Timestamp)
		if err != nil {
			log.Errorf("Repo log error %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Consume error %v", err)
	}
}
