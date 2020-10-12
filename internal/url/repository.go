package url

import (
	"fmt"
	"time"
	"turl/internal/cache"
	"turl/internal/db"
)

type Repository interface {
	Add(string, func(int64) (string, error)) (string, error)
	Get(int64) (string, error)
	GetCount(int64) (int, error)
	GetCountByTime(int64, time.Time) (int, error)
	Log(int64, time.Time) error
}

type URLRepository struct {
	db    *db.DB
	cache cache.Cache
}

func NewURLRepository(db *db.DB, cache cache.Cache) *URLRepository {
	return &URLRepository{db, cache}
}

func (r *URLRepository) Add(longUrl string, callback func(int64) (string, error)) (string, error) {
	var id int64
	err := r.db.QueryRow("INSERT INTO urls (url) VALUES ($1) RETURNING id", longUrl).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert url err: %v", err)
	}
	shortUrl, err := callback(id)
	if err != nil {
		return "", fmt.Errorf("url callback err: %v", err)
	}
	r.cache.Set(id, longUrl)

	return shortUrl, nil
}

func (r *URLRepository) Get(id int64) (string, error) {
	if longUrl, ok := r.cache.Get(id); ok {
		return longUrl, nil
	}
	var longUrl string
	err := r.db.QueryRow("SELECT url FROM urls WHERE id = $1", id).Scan(&longUrl)
	if err != nil {
		return "", fmt.Errorf("select url err: %v", err)
	}
	r.cache.Set(id, longUrl)

	return longUrl, nil
}

func (r *URLRepository) GetCount(id int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(id) FROM logs WHERE url_id = $1", id).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("select count err: %v", err)
	}

	return count, nil
}

func (r *URLRepository) GetCountByTime(id int64, from time.Time) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(id) FROM logs WHERE url_id = $1 AND datetime >= $2", id, from).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("select count by time err: %v", err)
	}

	return count, nil
}

func (r *URLRepository) Log(id int64, time time.Time) error {
	_, err := r.db.Exec("INSERT INTO logs (url_id, datetime) VALUES ($1, $2)", id, time)

	return err
}
