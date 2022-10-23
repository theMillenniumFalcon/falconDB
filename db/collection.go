package db

import (
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
)

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

func (cm *ConcurrentMap) GetRandomShard() *ConcurrentMapShared {
	return cm.Shared[rand.Intn(len(cm.Shared))]
}

func NewCollection(name string, files []*os.File, indexes []*Index, sd map[string]int) *Collection {
	return &Collection{name, NewConcurrentMap(files), indexes, sd, sync.Mutex{}, 0}
}

func (c *Collection) Size() uint64 {
	return atomic.LoadUint64(&c.ObjectsCounter)
}
