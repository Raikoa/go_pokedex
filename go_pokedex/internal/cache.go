package pokecache
import (
	"sync"
	"time"
	"fmt"
)
type Cache struct{
	Cargo   map[string]cacheEntry
	mu   sync.Mutex
	interval time.Duration
}


type cacheEntry struct{
	createdAt   time.Time
	val      []byte
}


func NewCache(inter time.Duration) *Cache{
	Newcache :=  &Cache{
		Cargo:    make(map[string]cacheEntry),
		interval:   inter,
	}
	
	go Newcache.reaploop(inter)

	return Newcache

}

func (c *Cache) Add(Key string, value []byte){
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheToAdd := cacheEntry{
		createdAt:  time.Now(),
		val:		value,
	}
	if _,exists := c.Cargo[Key]; !exists{
		c.Cargo[Key] = cacheToAdd
	}else{
		fmt.Print("already in cache")
		return
	}
}


func (c *Cache) Get(Key string) ([]byte, bool){
	c.mu.Lock()
	defer c.mu.Unlock()
	if _,exists := c.Cargo[Key]; !exists{
		return nil, false
	}else{
		return c.Cargo[Key].val, true
	}
}


func (c *Cache) reaploop(inter time.Duration){
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C{
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.Cargo{
			if now.Sub(entry.createdAt) > c.interval{
				fmt.Printf("removing stale entry")
				delete(c.Cargo, key)
			}
		}
		c.mu.Unlock()
	}
}


