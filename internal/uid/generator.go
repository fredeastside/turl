package uid

import (
	"fmt"
	"github.com/speps/go-hashids"
	"strconv"
)

type Generator interface {
	Encode(n int64) (string, error)
	Decode(s string) (int64, error)
}

type UIDGenerator struct {
	h *hashids.HashID
}

func NewUIDGenerator(salt, minLength string) (*UIDGenerator, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	var err error
	hd.MinLength, err = strconv.Atoi(minLength)
	if err != nil {
		return nil, fmt.Errorf("invalid min len UID generator error: %v", err)
	}
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, fmt.Errorf("init UID generator error: %v", err)
	}
	g := &UIDGenerator{h}

	return g, nil
}

func (g *UIDGenerator) Encode(n int64) (string, error) {
	s, err := g.h.EncodeInt64([]int64{n})
	if err != nil {
		return "", fmt.Errorf("encode UID error: %v", err)
	}

	return s, nil
}

func (g *UIDGenerator) Decode(s string) (int64, error) {
	n, err := g.h.DecodeInt64WithError(s)
	if err != nil {
		return 0, fmt.Errorf("decode UID error: %v", err)
	}

	return n[0], nil
}
