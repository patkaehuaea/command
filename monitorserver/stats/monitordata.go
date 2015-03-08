package stats

import (
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

type Sample struct {
    Time  time.Time
    Value int
}

type sequence map[string][]Sample

type MonitorData struct {
    sync.RWMutex
    data map[string]sequence
}

func (md *MonitorData) Add(target string, counter string, sample Sample) {
    md.Lock()
    md.data[target][counter] = append(md.data[target][counter], sample)
    md.Unlock()
}

func (md *MonitorData) Copy() map[string]sequence {
    copy := make(map[string]sequence)
    md.Lock()
    for url, seq := range md.data {
        copy[url] = make(sequence)
        for counter, slice := range seq {
            copy[url][counter] = slice
        }
    }
    md.Unlock()
    return copy
}

func New(keys []string) *MonitorData {
    temp := make(map[string]sequence)
    for _, key := range keys {
        temp[key] = make(sequence)
    }
    return &MonitorData{data: temp}
}

func (md *MonitorData) Print() {
    copy := md.Copy()
    json, _ := json.MarshalIndent(copy, "", " ")
    fmt.Println(string(json))
    return
}