// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package translator

type Translator interface {
	Translate(text string, languages []string) *TranslationContainer
	AddParticular(p particular.IParticularBackend)
}

