package imagecache

import (
	"image"
	"math/rand"
	"sync"
	"time"
)

type ImageCache struct {
	storage map[string]image.Image
	mux     sync.RWMutex
}

func NewImageCache() *ImageCache {
	return &ImageCache{storage: make(map[string]image.Image)}
}

func (ic *ImageCache) Get(key string) image.Image {
	ic.mux.RLock()
	defer ic.mux.RUnlock()
	return ic.storage[key]
}

func (ic *ImageCache) GetRandom() (imageID string, image image.Image) {
	ic.mux.RLock()
	defer ic.mux.RUnlock()

	rand.Seed(time.Now().Unix())
	randomIndex := rand.Intn(len(ic.storage))

	//TODO: make more effective
	for i, value := range ic.storage {
		if randomIndex == 0 {
			return i, value
		}
		randomIndex--
	}
	return "", nil
}

func (ic *ImageCache) Set(key string, image image.Image) {
	ic.mux.Lock()
	defer ic.mux.Unlock()
	ic.storage[key] = image
	return
}
