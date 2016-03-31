package cache_test

import (
	"log"
	"testing"

	"github.com/zsxm/scgo/data/cache"
)

func TestCacheSet(t *testing.T) {

	cache.Set("key1", "field1")
	val, _ := cache.Get("key1")
	log.Println(val)
}

func init() {
	cache.Conf = &cache.Config{
		Address:  "10.100.130.54:6379",
		Password: "foobared",
	}
	cache.Init(*cache.Conf)
}
