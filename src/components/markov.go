// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package components

import (
	"strings"
	"sync"
)

//http://www.wikiwand.com/en/Markov_chain
//https://github.com/simon-weber/gomarkov
//https://github.com/gilesp/markov/blob/master/markov.go

// Key
type IKey interface {
	String() string
	PushShift(s string)
}

func NewKey(size int) IKey {
	return &Key{words: make([]string, 0, size)}
}
func NewKeyFromSlice(s []string) IKey {
	return &Key{words: s}
}

type Key struct {
	words []string
}

func (k *Key) String() string {
	return strings.Join(k.words, " ")
}
func (k *Key) PushShift(word string) {
	copy(k.words, k.words[1:])
	k.words[len(k.words)-1] = word
}

// Suffixes
type ISuffixes interface {
	Add(s string)
	Occurrences(s string) float64
}

// Doubled space, @todo rewrite to trie
type Suffixes struct {
	occurrences map[string] float64
	suffixes []string
	sync.RWMutex
}
func NewSuffixes() ISuffixes {
	return &Suffixes{occurrences: make(map[string]float64)}
}

func (s *Suffixes) Add(word string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.occurrences[word]; ok {
		s.occurrences[word]++
	} else {
		s.suffixes = append(s.suffixes, word)
		s.occurrences[word] = 1
	}
}

func (s *Suffixes) Occurrences(word string) float64 {
	s.RLock()
	defer s.RUnlock()
	if occ, ok := s.occurrences[word]; ok {
		return occ
	}
	return 0
}

//Chain
type IChain interface {
	Add(s string)
	Occurrences(s string) float64
}

func NewChain(n int) IChain {
	return &Chain{grams: make(map[string]ISuffixes), n: n}
}

type Chain struct {
	grams map[string]ISuffixes
	n     int
}

func (c *Chain) Add(s string) {
	words := strings.Fields(s)
	if len(words) < c.n {
		return
	}

	// create key like [first second]
	key := NewKeyFromSlice(words[:c.n])

	for i := c.n; i < len(words); i++ {
		if _, ok := c.grams[key.String()]; !ok {
			c.grams[key.String()] = NewSuffixes()
		}
		c.grams[key.String()].Add(words[i])
		//c.grams[key.String()] = append(c.grams[key.String()], words[i])
		// update key to [second third]
		key.PushShift(words[i])
	}
}

func (c *Chain) Occurrences(s string) (res float64) {
	words := strings.Fields(s)
	if len(words) < c.n {
		return
	}

	key := NewKeyFromSlice(words[:c.n])

	for i := c.n; i < len(words); i++ {
		if _, ok := c.grams[key.String()]; ok {
			res += c.grams[key.String()].Occurrences(words[i])
		}
		key.PushShift(words[i])
	}
	return
}
