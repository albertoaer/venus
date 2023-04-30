package utils

import (
	"io"
	"math/rand"

	"github.com/oklog/ulid"
)

type IdGenerator interface {
	NextId() string
}

type ulidGenerator struct {
	entropy io.Reader
}

func NewUlidIdGenerator() IdGenerator {
	return &ulidGenerator{ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)}
}

func (gen *ulidGenerator) NextId() string {
	return ulid.MustNew(ulid.Now(), gen.entropy).String()
}
