package keylock

import (
	"sync"
)

func NewKeyLock() *Lock {
	return &Lock{lockMap: make(map[string]chan struct{})}
}

type Lock struct {
	lockMap map[string]chan struct{}
	mux     sync.RWMutex
}

func (l *Lock) GetLock(key string) chan struct{} {
	l.mux.Lock()
	defer l.mux.Unlock()
	keyLock := l.lockMap[key]

	if keyLock == nil {
		keyLock = make(chan struct{}, 1)

		l.lockMap[key] = keyLock

		return keyLock
	}

	return keyLock
}
