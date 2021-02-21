// requestlimiter limits amount of concurrent requests for any given url
package requestlimiter

import (
	"fmt"
	"sync"
)

var ErrQuotaReached = fmt.Errorf("request quota reached")
var ErrLockAlreadyAcquired = fmt.Errorf("lock already acquired")

type RequestLimiter struct {
	keyLock map[string]chan struct{}

	requestChan chan struct{}

	mux sync.RWMutex
}

func NewRequestLimiter(requestQuota int) *RequestLimiter {
	requestLimiter := RequestLimiter{keyLock: make(map[string]chan struct{})}

	if requestQuota > 0 {
		requestLimiter.requestChan = make(chan struct{}, requestQuota)
		for i := 0; i < requestQuota; i++ {
			requestLimiter.requestChan <- struct{}{}
		}
	}

	return &requestLimiter
}

// AcquireLock tries to acquire a lock for the given key. If succeed returns nil error.
// Call FreeLock after successfully acquiring the lock to avoid dead locks.
func (r *RequestLimiter) AcquireLock(key string) error {
	r.mux.RLock()

	keyLock := r.keyLock[key]

	r.mux.RUnlock()

	if keyLock != nil {
		<-keyLock
		return ErrLockAlreadyAcquired
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	if r.requestChan != nil {
		select {
		case <-r.requestChan:
		default:
			return ErrQuotaReached
		}
	}

	r.keyLock[key] = make(chan struct{})

	return nil
}

// FreeLock frees the lock acquired by AcquireLock
func (r *RequestLimiter) FreeLock(key string) {
	close(r.keyLock[key])
	if r.requestChan != nil {
		r.requestChan <- struct{}{}
	}
}
