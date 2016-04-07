// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package translator

import (
	"github.com/urakozz/go-emoji"
	"strconv"
	"regexp"
	"strings"
)

type EmojiProcessor struct {
	container map[string]string
}

var parser = emoji.NewEmojiParser()
var reviver = regexp.MustCompile(`\_\s?\$\s?\d+\_`)

func NewEmojiProcessor() *EmojiProcessor{
	return &EmojiProcessor{make(map[string]string)}
}

func (e *EmojiProcessor) Process(text string) string {
	i := -1
	return parser.ReplaceAllStringFunc(text, func(s string) string {
		i++
		key := " _$"+strconv.Itoa(i)+"_"
		e.container[key] = s
		return key
	})
}

func (e *EmojiProcessor) Restore(text string) string {
	return reviver.ReplaceAllStringFunc(text, func(s string) string {
		s = strings.Replace(s, " ", "", -1)
		return e.container[s]
	})
}

