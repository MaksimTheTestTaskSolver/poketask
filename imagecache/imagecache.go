package imagecache

import (
	"image"
	"sync"
)

type ImageCache struct {
	storage map[string]image.Image
	mux sync.RWMutex
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
	//TODO: return random value
	for i, value := range ic.storage {
		return i, value
	}
	return "", nil
}

func (ic *ImageCache) Set(key string, image image.Image) {
	ic.mux.Lock()
	defer ic.mux.Unlock()
	ic.storage[key] = image
	return
}
