package url

import (
	"fmt"
	"time"
	"turl/internal/uid"
)

type Shortener interface {
	Encode(string) (string, error)
	Decode(string) (string, error)
	GetDailyCount(string) (int, error)
	GetWeeklyCount(string) (int, error)
	GetCount(string) (int, error)
}

type URLShortener struct {
	repo     Repository
	producer Producer
	g        uid.Generator
}

func NewURLShortener(repo Repository, producer Producer, g uid.Generator) *URLShortener {
	return &URLShortener{repo, producer, g}
}

func (sh *URLShortener) Encode(s string) (string, error) {
	shortUrl, err := sh.repo.Add(s, func(id int64) (string, error) {
		return sh.g.Encode(id)
	})
	if err != nil {
		return "", fmt.Errorf("shortener encode err: %v", err)
	}

	return shortUrl, nil
}

func (sh *URLShortener) Decode(s string) (string, error) {
	id, err := sh.g.Decode(s)
	if err != nil {
		return "", fmt.Errorf("shortener decode err: %v", err)
	}
	longUrl, err := sh.repo.Get(id)
	if err != nil {
		return "", fmt.Errorf("shortener get err: %v", err)
	}
	err = sh.producer.Produce(id, time.Now().UTC())
	if err != nil {
		return "", fmt.Errorf("shortener save to queue err: %v", err)
	}

	return longUrl, nil
}

func (sh *URLShortener) GetDailyCount(s string) (int, error) {
	count, err := sh.getCount(s, func(id int64) (int, error) {
		return sh.repo.GetCountByTime(id, time.Now().AddDate(0, 0, -1))
	})
	if err != nil {
		return 0, fmt.Errorf("shortener daily count err: %v", err)
	}

	return count, nil
}

func (sh *URLShortener) GetWeeklyCount(s string) (int, error) {
	count, err := sh.getCount(s, func(id int64) (int, error) {
		return sh.repo.GetCountByTime(id, time.Now().AddDate(0, 0, -7))
	})
	if err != nil {
		return 0, fmt.Errorf("shortener weekly count err: %v", err)
	}

	return count, nil
}

func (sh *URLShortener) GetCount(s string) (int, error) {
	count, err := sh.getCount(s, func(id int64) (int, error) {
		return sh.repo.GetCount(id)
	})
	if err != nil {
		return 0, fmt.Errorf("shortener all count err: %v", err)
	}

	return count, nil
}

func (sh *URLShortener) getCount(s string, callback func(int64) (int, error)) (int, error) {
	id, err := sh.g.Decode(s)
	if err != nil {
		return 0, fmt.Errorf("shortener decode err: %v", err)
	}

	return callback(id)
}
