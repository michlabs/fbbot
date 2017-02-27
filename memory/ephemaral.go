package memory

import (
	"sync"
)

type ephemeralMemory struct {
	mutex   *sync.Mutex
	mapping map[string]Store
}

func newEphemeralMemory() *ephemeralMemory {
	return &ephemeralMemory{
		mutex:   &sync.Mutex{},
		mapping: make(map[string]Store),
	}
}

func (em ephemeralMemory) For(id string) Store {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	es, ok := em.mapping[id]
	if !ok {
		s := newEphemeralStore()
		em.mapping[id], es = s, s
	}
	return es
}

func (em ephemeralMemory) Delete(id string) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	delete(em.mapping, id)
}

// ephemeralStore is a memory that stores data in RAM
// Just use it for development
type ephemeralStore struct {
	mutex *sync.Mutex
	store map[string]string
}

func newEphemeralStore() Store {
	return &ephemeralStore{
		mutex: &sync.Mutex{},
		store: make(map[string]string),
	}
}

func (es ephemeralStore) Set(key string, value string) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	es.store[key] = value
}

func (es ephemeralStore) Get(key string) string {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	return es.store[key]
}

func (es ephemeralStore) Delete(key string) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	delete(es.store, key)
}
