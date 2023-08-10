package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

var MemCache *memcache.Client

func Init() {
	MemCache = memcache.New("127.0.0.1:11211")
}
