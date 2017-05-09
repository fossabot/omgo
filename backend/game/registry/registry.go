package registry

import "sync"

// Registry is a in memory map for recording user information
type Registry struct {
	sync.RWMutex
	records map[uint64]interface{} // id -> v
}

var (
	defaultRegistry Registry
)

func init() {
	defaultRegistry.init()
}

func (r *Registry) init() {
	r.records = make(map[uint64]interface{})
}

// Register add an user to record and use user's ID as key
func (r *Registry) Register(id uint64, v interface{}) {
	r.Lock()
	r.records[id] = v
	r.Unlock()
}

// Unregister removes user with specific user ID from record
func (r *Registry) Unregister(id uint64, v interface{}) {
	r.Lock()
	if oldv, ok := r.records[id]; ok {
		if oldv == v {
			delete(r.records, id)
		}
	}
	r.Unlock()
}

// Query user with ID
func (r *Registry) Query(id uint64) (v interface{}) {
	r.RLock()
	v = r.records[id]
	r.RUnlock()
	return
}

// Count returns amount of online users
func (r *Registry) Count() (count int) {
	r.RLock()
	count = len(r.records)
	r.RUnlock()
	return
}

// Register an user to default registry
func Register(id uint64, v interface{}) {
	defaultRegistry.Register(id, v)
}

// Unregister an user from default registry
func Unregister(id uint64, v interface{}) {
	defaultRegistry.Unregister(id, v)
}

// Query user in default registry
func Query(id uint64) interface{} {
	return defaultRegistry.Query(id)
}

// Count users in default registry
func Count() int {
	return defaultRegistry.Count()
}
