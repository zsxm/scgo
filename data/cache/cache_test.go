package cache_test

import (
	"log"
	"testing"
	"weixin/source/pubnum/entity"

	"github.com/zsxm/scgo/data/cache"
)

func TestCacheSet(t *testing.T) {
	mp := make(map[string]string)
	mp["user"] = "zhagsan"
	mp["age"] = "24"
	cache.HSetMap("users", mp)
	cache.HGetMap("users")
	ent := entity.NewPubnum()
	ent.SetId("id")
	ent.SetName("张三")
	cache.HSetEntity("pubnum", ent)
	ent = entity.NewPubnum()
	cache.HGetEntity("pubnum", ent)
	log.Println(ent)
}

func init() {
	cache.Conf = &cache.Config{
		Address:  "127.0.0.1:6379",
		Password: "foobared",
	}
	cache.Init(*cache.Conf)
}
