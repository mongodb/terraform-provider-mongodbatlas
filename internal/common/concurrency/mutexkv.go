package concurrency

import (
	"log"
	"sync"
)

// Copied from Hashicorp: https://developer.hashicorp.com/terraform/plugin/sdkv2/guides/v2-upgrade-guide#removal-of-helper-mutexkv-package
// "Providers that need the functionality provided by the helper/mutexkv package are encouraged to copy the types and functions
// it provided into their own codebase, provided here under a public domain license"

// MutexKV is a simple key/value store for arbitrary mutexes. It can be used to
// serialize changes across arbitrary collaborators that share knowledge of the
// keys they must serialize on.
type MutexKV struct {
	store map[string]*sync.Mutex
	lock  sync.Mutex
}

// Locks the mutex for the given key. Caller is responsible for calling Unlock
// for the same key
func (m *MutexKV) Lock(key string) {
	log.Printf("[DEBUG] Locking %q", key)
	m.get(key).Lock()
	log.Printf("[DEBUG] Locked %q", key)
}

// Unlock the mutex for the given key. Caller must have called Lock for the same key first
func (m *MutexKV) Unlock(key string) {
	log.Printf("[DEBUG] Unlocking %q", key)
	m.get(key).Unlock()
	log.Printf("[DEBUG] Unlocked %q", key)
}

// Returns a mutex for the given key, no guarantee of its lock status
func (m *MutexKV) get(key string) *sync.Mutex {
	m.lock.Lock()
	defer m.lock.Unlock()
	mutex, ok := m.store[key]
	if !ok {
		mutex = &sync.Mutex{}
		m.store[key] = mutex
	}
	return mutex
}

// Returns a properly initialized MutexKV
func NewMutexKV() *MutexKV {
	return &MutexKV{
		store: make(map[string]*sync.Mutex),
	}
}
