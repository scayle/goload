package random

import (
	lfring "github.com/LENSHOOD/go-lock-free-ring-buffer"
	"golang.org/x/sync/singleflight"
	"math/rand"
)

type Sampler[T any] struct {
	Values []T
	rb     lfring.RingBuffer[T]
	sfg    *singleflight.Group
}

func NewSampler[T any](values []T) *Sampler[T] {
	return &Sampler[T]{
		Values: values,
		rb:     lfring.New[T](lfring.NodeBased, uint64(len(values)*2)),
		sfg:    &singleflight.Group{},
	}
}

// Get returns n 'random' items
// because the values gets preshuffled in a ring buffer there is no guaranty that there are no duplicate items
func (s *Sampler[T]) Get(n int) []T {
	values := make([]T, 0, n)
	for len(values) < n {
		val, ok := s.rb.Poll()
		if !ok {
			_, _, _ = s.sfg.Do("fillBuffer", func() (interface{}, error) {
				s.fillBuffer()
				return nil, nil
			})
			continue
		}
		values = append(values, val)
	}
	return values
}

func (s *Sampler[T]) fillBuffer() {
	rand.Shuffle(len(s.Values), func(i, j int) { s.Values[i], s.Values[j] = s.Values[j], s.Values[i] })
	for _, value := range s.Values {
		s.rb.Offer(value)
	}
}
