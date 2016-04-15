// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package processing

import "unicode"

type ProcessingSegments struct {

}

func (p *ProcessingSegments) some(){
	unicode.IsGraphic("s")
}
