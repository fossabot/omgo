package registry

import "sync"

type Registry struct {
	sync.RWMutex
	entries map[uint64]interface{} // usn -> v
}

var (
	defaultRegistry Registry
)

func init() {
	defaultRegistry.init()
}

func (r *Registry) init() {
	r.entries = make(map[uint64]interface{})
}

// Register adds a new entry
func (r *Registry) Register(usn uint64, v interface{}) {
	r.Lock()
	r.entries[usn] = v
	r.Unlock()
}

// Unregister removes an entry from entries
func (r *Registry) Unregister(usn uint64, v interface{}) {
	r.Lock()
	if oldv, ok := r.entries[usn]; ok {
		if oldv == v {
			delete(r.entries, usn)
		}
	}
	r.Unlock()
}

// Query an entry
func (r *Registry) Query(usn uint64) interface{} {
	r.RLock()
	defer r.RUnlock()
	return r.entries[usn]
}

// Count entries
func (r *Registry) Count() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.entries)
}

// Register API
func Register(usn uint64, v interface{}) {
	defaultRegistry.Register(usn, v)
}

// Unregister API
func Unregister(usn uint64, v interface{}) {
	defaultRegistry.Unregister(usn, v)
}

// Query API
func Query(usn uint64) interface{} {
	return defaultRegistry.Query(usn)
}

// Count API
func Count() int {
	return defaultRegistry.Count()
}
