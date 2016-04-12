// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package backend_particular

type Meaning struct {
	Dictionary string
	Text       string
	All        []string
}

func (m *Meaning) GetText() string {
	return m.Text
}
func (m *Meaning) GetAll() []string {
	return m.All
}

func (m *Meaning) GetDictionary() string {
	return m.Dictionary
}


var _ IParticularMeaning = (*Meaning)(nil)

