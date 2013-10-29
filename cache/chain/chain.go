// Chained cache
package chain

import (
	"github.com/pierrre/imageserver"
)

type ChainCache []imageserver.Cache

// Get an image from caches in sequential order
//
// If an image is found, previous caches are filled
func (cache ChainCache) Get(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
	for i, c := range cache {
		image, err := c.Get(key, parameters)

		if err == nil {
			if i > 0 {
				cache.setCaches(key, image, parameters, i)
			}
			return image, nil
		}

		if _, ok := err.(*imageserver.CacheMissError); !ok {
			return nil, err
		}
	}

	return nil, imageserver.NewCacheMissError(key, cache, nil)
}

func (cache ChainCache) setCaches(key string, image *imageserver.Image, parameters imageserver.Parameters, indexLimit int) {
	for i := 0; i < indexLimit; i++ {
		go func(i int) {
			cache[i].Set(key, image, parameters)
		}(i)
	}
}

// Set the image to all caches
func (cache ChainCache) Set(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
	for _, c := range cache {
		go func(c imageserver.Cache) {
			c.Set(key, image, parameters)
		}(c)
	}
	return nil
}
