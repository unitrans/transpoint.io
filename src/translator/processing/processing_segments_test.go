// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package processing

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestSentences(t *testing.T) {
	segs := Segments.Split("hello. How are you?\n I was ... ")
	correct := []string{"hello", ". ", "How are you", "?\n ", "I was ", "... "}

	for i, s := range segs{
		assert.Equal(t, correct[i], s.Text)
	}
	assert.Equal(t, 6, len(segs))
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
func TestEndsOnDot(t *testing.T) {
	segs := Segments.Split("Sentence boundary detection is hard. Nasennebenhöhlenentzündung.")

	assert.Equal(t, 4, len(segs))
}
func TestSimpleEmoji(t *testing.T) {
	str := `like 💩and 🍦
Skin tone options 🎅🏻🎅🏼🎅🏽🎅🏾🎅🏿. iOS brings skin tone options to existing emoji such as Santa Claus
😳The regular expression
`
	strList := strings.FieldsFunc(str, func(r rune) bool {
		return r == '\n'
	})
	correctLen := []int{4,3,2}
	for i, sent := range strList{
		segs := Segments.Split(sent)
		assert.Equal(t, correctLen[i], len(segs))
		//for j, seg := range segs{
		//
		//	fmt.Println(i, j, seg.Text)
		//}
	}
}

