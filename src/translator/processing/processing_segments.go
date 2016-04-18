// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package processing




type SegmentType int

const (
	SegmentText SegmentType = iota
	SegmentSeparator
)

type Segment struct {
	Type SegmentType
	Text string
}


var Segments = &ProcessingSegments{}

//@todo write proper automata
type ProcessingSegments struct {

}


func (p *ProcessingSegments) Split(s string) []*Segment {
	//unicode.IsGraphic([]rune("s"))
	return p.FieldsFunc(s, p.CombineFuncs(p.IsNewLine, p.IsDot, p.IsPunctuation))
}

func (p *ProcessingSegments) CombineFuncs(fns... func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		for _, fn := range fns {
			if fn(r) {
				return true
			}
		}
		return false
	}
}

func (p *ProcessingSegments) IsNewLine(r rune) bool {
	switch r {
	case '\n', 0x85:
		return true
	}
	return false

}

func (p *ProcessingSegments) IsDot(r rune) bool {
	switch r {
	case '.', 0x2024, 0x2025, 0x2026:
		return true
	}
	return false

}

func (p *ProcessingSegments) IsPunctuation(r rune) bool {
	switch r {
	case '!', '?', 0x2047, 0x2048, 0x2049 :
		return true
	}
	return false

}

func (p *ProcessingSegments) FieldsFunc(s string, f func(rune) bool) []*Segment {
	// First count the fields.
	n := 0
	inField := false
	for _, rune := range s {
		wasInField := inField
		inField = !f(rune)
		if inField && !wasInField {
			n++
		}
	}
	if n == 0 {
		return []*Segment{}
	}

	// Now create them.
	a := make([]*Segment, n*2-1)
	na := 0
	//fieldStart := -1 // Set to -1 when looking for start of field.
	for i, rune := range s {
		if i == 0 {
			a[na] = &Segment{Text:string(rune),Type:SegmentText}
			if f(rune) {
				a[na].Type = SegmentSeparator
			}
		} else {
			t := SegmentText
			if f(rune) {
				t = SegmentSeparator
			}
			if a[na].Type == t {
				a[na].Text += string(rune)
			} else {
				na++
				a[na] = &Segment{Text:string(rune),Type:t}
			}

		}

	}

	return a
}
