package imagecache

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/MaksimTheTestTaskSolver/poketask/imagecache/keylock"
	"github.com/MaksimTheTestTaskSolver/poketask/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var ErrTimeout = fmt.Errorf("timeout")

func NewImageCache() *ImageCache {
	return &ImageCache{
		storage: make(map[string]*model.Image),
		keyLock: keylock.NewKeyLock(),
	}
}

type ImageCache struct {
	storage map[string]*model.Image
	keyLock *keylock.Lock
	mux     sync.RWMutex
}

func (ic *ImageCache) Get(key string) *model.Image {
	ic.mux.RLock()
	defer ic.mux.RUnlock()
	return ic.storage[key]
}

func (ic *ImageCache) GetRandom() (image *model.Image) {
	ic.mux.RLock()
	defer ic.mux.RUnlock()

	cachedImagesLen := len(ic.storage)

	if cachedImagesLen == 0 {
		return nil
	}

	randomIndex := rand.Intn(cachedImagesLen)

	//TODO: make more effective
	for _, value := range ic.storage {
		if randomIndex == 0 {
			return value
		}
		randomIndex--
	}
	return nil
}

func (ic *ImageCache) Set(key string, image *model.Image) {
	ic.mux.Lock()
	defer ic.mux.Unlock()

	if storedImage := ic.storage[key]; storedImage != nil {
		storedImage = image
		return
	}

	ic.storage[key] = image

	return
}

// GetWithLock return the image if it exists in the cache or nil otherwise.
// If image does not exist it locks access to the image until unlock function is called
// If the lock already acquired it waits till lock is unlocked or timeouts after 10 seconds with ErrTimeout
//
// Always call the unlock function to prevent deadlocks:
//
// image, unlock, err := GetWithLock(...)
// if err != nil {...}
// defer unlock()
// if image == nil {
// 	...
// 	Set(...)
// }
//
func (ic *ImageCache) GetWithLock(key string) (image *model.Image, unlock func(), err error) {
	ic.mux.RLock()
	image = ic.storage[key]
	ic.mux.RUnlock()

	if image == nil {
		keyLock := ic.keyLock.GetLock(key)

		select {
		case keyLock<- struct{}{}:
		case <-time.After(10 * time.Second):
			return nil, func() {}, ErrTimeout
		}

		// check if value was set while we were waiting for the lock
		if image := ic.storage[key]; image != nil {
			<-keyLock
			return image, func() {}, nil
		}

		return nil, func() { <-keyLock }, nil
	}

	return image, func() {}, nil
}
