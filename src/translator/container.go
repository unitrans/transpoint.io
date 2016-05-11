// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package translator

import (
	"github.com/unitrans/unitrans/src/translator/backend_particular"
	"time"
	"sort"
)

// TranslationBag hashmap
type TranslationBag map[string]string

// RawTransData struct
type RawTranslationData struct {
	Source      string
	Lang        string
	Name        string
	Original    string
	Translation string
	Time        time.Duration
	Score       float64
}

type RawParticularData struct {
	Original      string
	Source        string
	Lang          string
	Name          string
	ParticularBag []*ParticularItem
	Time          time.Duration
}

type ParticularItem struct {
	Url          string
	Order        int
	Original     string
	Time         time.Duration
	Translations []backend_particular.IParticularMeaning
}

// TranslationContainer struct
type TranslationContainer struct {
	Langs             []string
	Translations      TranslationBag
	Original          string
	Source            string
	//Meta              map[string]interface{}
	//RawTranslations   map[string]TranslationBag
	RawTransData      map[string]map[string]*RawTranslationData
	RawSegmentsData   []map[string]map[string]*RawTranslationData //`json:"-"`
}


func NewSegmentsSorter(m map[string]*RawTranslationData) *segmentSorter {
	list := make([]*RawTranslationData, 0, len(m))
	for _, v := range m{
		list = append(list, v)
	}
	return &segmentSorter{segments:list}
}

type segmentSorter struct {
	segments []*RawTranslationData
}

// Len is part of sort.Interface.
func (s *segmentSorter) Len() int {
	return len(s.segments)
}

// Swap is part of sort.Interface.
func (s *segmentSorter) Swap(i, j int) {
	s.segments[i], s.segments[j] = s.segments[j], s.segments[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *segmentSorter) Less(i, j int) bool {
	return s.segments[i].Score > s.segments[j].Score
}

func (s *segmentSorter) Max() *RawTranslationData {
	if s.Len() == 0 {
		return nil
	}
	sort.Sort(s)
	return s.segments[0]
}


