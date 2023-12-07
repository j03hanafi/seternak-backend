package id

import (
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"sync"
	"time"
)

var entropyPool = sync.Pool{
	New: func() any {
		rnd := rand.Reader
		entropy := ulid.Monotonic(rnd, 0)
		return entropy
	},
}

func New() ulid.ULID {
	entropy := entropyPool.Get().(ulid.MonotonicReader)
	defer entropyPool.Put(entropy)

	timestamp := ulid.Timestamp(time.Now())
	return ulid.MustNew(timestamp, entropy)
}
