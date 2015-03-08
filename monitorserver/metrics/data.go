//  Copyright (C) Pat Kaehuaea - All Rights Reserved
//  Unauthorized copying of this file, via any medium is strictly prohibited
//  Proprietary and confidential
//  Written by Pat Kaehuaea, March 2015
//
// Stores collected monitoring data for monitorserver. Provides
// interface to caller of dictinary whose primary key is a target
// host, and whose value is a map of counters to type sequence.
// The Package implements an Add(), Copy(), New(), and Print()
// function. Calling Print() will write a JSON encoded string
// of the sequences to the screen.

package metrics

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Intended for storage of counter value and
// the time it was collected.
type Sample struct {
	Time  time.Time
	Value int
}

// A map whose key is a counter value e.g. "500s"
// and value is a slice of Sample structs.
type sequence map[string][]Sample

type Data struct {
	sync.RWMutex
	sequences map[string]sequence
}

// Add a new Sample to the Data for target and counter.
func (d *Data) Add(target string, counter string, sample Sample) {
	d.Lock()
	d.sequences[target][counter] = append(d.sequences[target][counter], sample)
	d.Unlock()
}

// Returns map of targets as strings to their collected
// metrics (sequences).
func (d *Data) Copy() map[string]sequence {
	copy := make(map[string]sequence)
	d.Lock()
	for url, seq := range d.sequences {
		copy[url] = make(sequence)
		for counter, slice := range seq {
			copy[url][counter] = slice
		}
	}
	d.Unlock()
	return copy
}

// Returns pointer to new Data struct whose
// primary keys are passed via keys parameter.
func New(keys []string) *Data {
	temp := make(map[string]sequence)
	for _, key := range keys {
		temp[key] = make(sequence)
	}
	return &Data{sequences: temp}
}

// Prints the JSON representation of sequences
// data element of Data struct.
func (d *Data) Print() {
	copy := d.Copy()
	json, _ := json.MarshalIndent(copy, "", " ")
	fmt.Println(string(json))
	return
}
