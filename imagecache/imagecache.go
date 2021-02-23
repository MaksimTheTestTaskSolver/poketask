package imagecache

import (
	"math/rand"
	"sync"
	"time"

	"github.com/MaksimTheTestTaskSolver/poketask/model"
)

type ImageCache struct {
	storage map[string]*model.Image
	mux     sync.RWMutex
}

func NewImageCache() *ImageCache {
	return &ImageCache{storage: make(map[string]*model.Image)}
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

	rand.Seed(time.Now().Unix())
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
	ic.storage[key] = image
	return
}
