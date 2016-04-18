// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package processing

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSentences(t *testing.T) {
	segs := Segments.Split("hello. How are you?\n I was ... ")
	correct := []string{"hello", ".", " How are you", "?\n", " I was ", "...", " "}

	for i, s := range segs{
		assert.Equal(t, correct[i], s.Text)
	}
	assert.Equal(t, 7, len(segs))
}
func TestSpace(t *testing.T) {
	segs := Segments.Split(" ")
	correct := []string{" "}

	for i, s := range segs{
		assert.Equal(t, correct[i], s.Text)
	}
	assert.Equal(t, 1, len(segs))
}
func TestEmpty(t *testing.T) {
	segs := Segments.Split("")

	assert.Equal(t, 0, len(segs))
}

