package syncjob

import (
	"sync"
)

var (
	l = sync.Mutex{}

	jobCache = make(map[string]int, 50)

	JobChannel = make(chan string)
)

func Get() map[string]int {
	return jobCache
}

func Set(key string, v int) {
	l.Lock()
	defer l.Unlock()
	jobCache[key] = v
}

func Clean(in map[string]int) {
	l.Lock()
	defer l.Unlock()
	for k := range in {
		delete(jobCache, k)
	}
}

func RangeSet(in map[string]int) {
	l.Lock()
	defer l.Unlock()
	for k, v := range in {
		jobCache[k] = v
	}
}
