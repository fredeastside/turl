package main

import (
	"os"
	"turl/internal/cache"
	"turl/internal/db"
	"turl/internal/server"
	"turl/internal/uid"
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

	g, err := uid.NewUIDGenerator(os.Getenv("UID_SALT"), os.Getenv("UID_MIN_LEN"))
	if err != nil {
		log.Fatalf("Init UID error %s", err.Error())
	}
	p, err := url.NewKafkaProducer(os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	if err != nil {
		log.Fatalf("Kafka conn error %s", err.Error())
	}
	defer p.Close()
	sh := url.NewURLShortener(url.NewURLRepository(conn, cache.NewInMemoryCache()), p, g)
	s := server.NewServer(sh)
	log.Fatal(s.Run(":" + os.Getenv("PORT")))
}
