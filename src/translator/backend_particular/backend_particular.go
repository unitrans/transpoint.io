// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package backend_particular

type IBackendParticular interface {
	TranslateWord(text string, language, to string) (data IBackendParticularResponse)
	GetName() string
}

type IBackendParticularResponse interface {
	GetMeanings() []IParticularMeaning
	GetUrl() string
}

type IParticularMeaning interface {
	GetDictionary() string
	GetText() string
	GetAll() []string
}
