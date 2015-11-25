// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package translation_middleware
import (
	"strconv"
	"regexp"
	"github.com/urakozz/go-emoji"
	"github.com/urakozz/transpoint.io/src/infrastrucrute/translator"
)

type EmojiProcessor struct {
	container map[string]string
}

var parser = emoji.NewEmojiParser()
var reviver = regexp.MustCompile(`\_\$\d+\_`)

func NewEmojiProcessor() *EmojiProcessor{
	return &EmojiProcessor{make(map[string]string)}
}

func (e *EmojiProcessor) Process(text string) string {
	i := -1
	return parser.ReplaceAllStringFunc(text, func(s string) string {
		i++
		key := "_$"+strconv.Itoa(i)+"_"
		e.container[key] = s
		return key
	})
}

func (e *EmojiProcessor) Restore(text string) string {
	return reviver.ReplaceAllStringFunc(text, func(s string) string {
		return e.container[s]
	})
}

type EmojiMiddleware struct {

}

func (m *EmojiMiddleware) MiddlewareFunc (h HandlerFunc) HandlerFunc {
	return func(c *translator.TranslationContainer) {
		processor := NewEmojiProcessor()
		c.Original = processor.Process(c.Original)

		h(c)

		for k, v := range c.Translations {
			c.Translations[k] = processor.Restore(v)
		}
	}

}