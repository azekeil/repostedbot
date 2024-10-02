package reposted

import "sync"

type SafeRWMutex struct {
	*sync.RWMutex
}

func (sm *SafeRWMutex) SRLock() {
	if sm.RWMutex == nil {
		sm.RWMutex = &sync.RWMutex{}
	}
	sm.RLock()
}

func (sm *SafeRWMutex) SLock() {
	if sm.RWMutex == nil {
		sm.RWMutex = &sync.RWMutex{}
	}
	sm.Lock()
}
