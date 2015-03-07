package sequence

import (
    "sync"
    "time"
)

type Sample struct {
	Time  time.Time
	Value int
}

type Sequence struct {
	sync.RWMutex
	sequences map[string][]Sample
}

func (s *Sequence) Add(counter string, sample Sample) {
    s.Lock()
    s.sequences[counter] = append(s.sequences[counter], sample)
    s.Unlock()
}

func (s *Sequence) Copy() (copy map[string][]Sample) {
    copy = make(map[string][]Sample)
    s.Lock()
    for counter, samples := range s.sequences {
        copy[counter] = samples
    }
    s.Unlock()
    return
}

func New() *Sequence {
    return &Sequence{sequences: make(map[string][]Sample)}
}


