// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package backend_full

type IBackendFullResponse interface {
	GetText() string
	GetSource() string
	GetLang() string
}

type IBackendFull interface {
	TranslateFull(text string, language string) (data IBackendFullResponse)
	GetName() string
}
