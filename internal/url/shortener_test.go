package url

import (
	"errors"
	"strconv"
	"testing"
	"time"
	"turl/internal/uid"
)

type mockGenerator struct {
}

func (m *mockGenerator) Encode(n int64) (string, error) {
	return strconv.Itoa(int(n)), nil
}

func (m *mockGenerator) Decode(s string) (int64, error) {
	n, _ := strconv.Atoi(s)

	return int64(n), nil
}

type mockProducer struct {
	repo Repository
}

func (m *mockProducer) Produce(id int64, t time.Time) error {
	m.repo.Log(id, t)

	return nil
}

type mockRepository struct {
	count int64
	urls  map[int64]string
	logs  map[int64]time.Time
	g     uid.Generator
}

func (m *mockRepository) Add(s string, f func(int64) (string, error)) (string, error) {
	m.count++
	m.urls[m.count] = s

	return m.g.Encode(m.count)
}

func (m *mockRepository) Get(id int64) (string, error) {
	if _, ok := m.urls[id]; !ok {
		return "", errors.New("error")
	}

	return m.urls[id], nil
}

func (m *mockRepository) GetCount(int64) (int, error) {
	return len(m.logs), nil
}

func (m *mockRepository) GetCountByTime(int64, time.Time) (int, error) {
	return len(m.logs), nil
}

func (m *mockRepository) Log(id int64, t time.Time) error {
	m.logs[id] = t

	return nil
}

func TestEncode(t *testing.T) {
	sh := getURLShortener()
	s := "12"
	v, err := sh.Encode(s)
	if err != nil {
		t.Fatalf("You received an error %v.", err)
	}
	if v != "1" {
		t.Fatalf("Expected 1 got %v.", v)
	}
}

func TestDecode(t *testing.T) {
	sh := getURLShortener()
	s := "12"
	v, _ := sh.Encode(s)
	v, err := sh.Decode(v)
	if err != nil {
		t.Fatalf("You received an error %v.", err)
	}
	if v != s {
		t.Fatalf("Expected %v got %v.", s, v)
	}
}

func TestGetCount(t *testing.T) {
	sh := getURLShortener()
	s := "12"
	count := 1
	v, _ := sh.Encode(s)
	v, _ = sh.Decode(v)
	c, err := sh.GetCount(v)
	if err != nil {
		t.Fatalf("You received an error %v.", err)
	}
	if c != count {
		t.Fatalf("Expected %v got %v.", count, c)
	}
}

func getURLShortener() *URLShortener {
	g := &mockGenerator{}
	repo := &mockRepository{urls: make(map[int64]string), logs: make(map[int64]time.Time), g: g}
	producer := &mockProducer{repo}

	return NewURLShortener(repo, producer, g)
}
