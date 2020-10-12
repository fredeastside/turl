package uid

import "testing"

func TestGenerator(t *testing.T) {
	g, err := NewUIDGenerator("test", "5")
	if err != nil {
		t.Error(err.Error())
	}
	v, _ := g.Encode(int64(1))
	id, _ := g.Decode(v)
	if id != int64(1) {
		t.Errorf("Expected %v got %v", int64(1), id)
	}
}
