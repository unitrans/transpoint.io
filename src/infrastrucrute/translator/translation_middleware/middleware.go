// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package translation_middleware
import "github.com/urakozz/transpoint.io/src/infrastrucrute/translator"

type HandlerFunc func(c *translator.TranslationContainer)

type ITranslationMiddleware interface{
	MiddlewareFunc(handler HandlerFunc) HandlerFunc
}
