// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package particular
import "net/http"

type IParticularBackend interface {
	TranslateOne(text string, language, to string) (data IParticularResponse)
	GetName() string
	SetClient(c *http.Client)
}

type IParticularResponse interface {
	GetMeanings() []IParticularMeaning
}

type Meaning struct {
	Dictationary string
	Text string
	All []string
}

type IParticularMeaning interface {
	GetDictationary() string
	GetText() string
	GetAll() []string
}

func (m *Meaning) GetText() string {
	return m.Text
}
func (m *Meaning) GetAll() []string {
	return m.All
}

func (m *Meaning) GetDictationary() string {
	return m.Dictationary
}

