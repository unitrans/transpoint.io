// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package processing

import "strings"

//import "unicode"

type ProcessingSegments struct {

}

func (p *ProcessingSegments) Split(s string) []strings {
	//unicode.IsGraphic([]rune("s"))
	return strings.FieldsFunc(s, p.CombineFuncs(IsNewLine))
}

func (p *ProcessingSegments) CombineFuncs(fns... func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		for fn := range fns {
			if fn(r) {
				return true
			}
		}
		return false
	}
}

func IsNewLine(r rune) bool {
	switch r {
	case '\n', 0x85:
		return true
	}
	return false

}
