package cache

import "testing"

func TestCache(t *testing.T) {
	id := int64(1)
	val := "test"
	c := NewInMemoryCache()
	c.Set(id, val)
	v, ok := c.Get(id)
	if !ok {
		t.Error("expected true got false")
	}
	if v != val {
		t.Error("expected true got false")
	}
	_, ok = c.Get(id + 1)
	if ok {
		t.Error("expected false got true")
	}
}
