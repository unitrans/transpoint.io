// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	//"fmt"
)

func TestAdd(t *testing.T) {
	m := NewChain(2)
	m.Add("we have to log data which has been read by which client")
	m.Add("by which kind of issue it have been logged out")
	assert.Equal(t, 17, len(m.(*Chain).grams))
	suffixes, ok := m.(*Chain).grams["we have"]
	suffixes2, ok2 := m.(*Chain).grams["have to"]
	suffixes3, ok3 := m.(*Chain).grams["by which"]
	assert.True(t, ok)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, suffixes.(*Suffixes).suffixes, []string{"to"})
	assert.Equal(t, suffixes.(*Suffixes).occurrences["to"], float64(1))
	assert.Equal(t, suffixes2.(*Suffixes).suffixes, []string{"log"})
	assert.Equal(t, suffixes2.(*Suffixes).occurrences["log"], float64(1))
	assert.Equal(t, suffixes3.(*Suffixes).suffixes, []string{"client", "kind"})
	assert.Equal(t, suffixes3.(*Suffixes).occurrences["client"], float64(1))
	assert.Equal(t, suffixes3.(*Suffixes).occurrences["kind"], float64(1))
	//fmt.Println(suffixes3.(*Suffixes))
}

func TestOcc(t *testing.T) {
	m := NewChain(2)
	m.Add("we have to log data which has been read by which client")
	m.Add("by which kind of issue it have been logged out")
	m.Add("controls the size of a core dump for a process running as a named user or group")
	m.Add("system call to create a new process")
	assert.Equal(t, float64(10), m.Occurrences("we have to log data which has been read by which client"))
	assert.Equal(t, float64(2), m.Occurrences("which has been changed by which client"))
}

