package db

import "sync"

type Index struct {
	Field  string
	Unique bool
}

type IndexData struct {
	Field string
	Data  string
}

type Collection struct {
	Name              string
	Map               *ConcurrentMap
	Indexes           []*Index
	ShardDestinations map[string]int
	mx                sync.Mutex
	ObjectsCounter    uint64
}

type Element struct {
	Id      string      `json:"id"`
	Payload interface{} `json:"payload"`
}
