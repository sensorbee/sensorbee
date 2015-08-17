package exp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

// CacheEntry has an cache information of a specific BQL statement.
type CacheEntry struct {
	NodeName  string    `json:"node_name,omitempty"`
	Hash      string    `json:"hash"`
	Stmt      string    `json:"stmt"`
	States    []string  `json:"states"`
	Timestamp time.Time `json:"timestamp"`
}

// Cache has information of cached results generated from each BQL statement.
type Cache struct {
	Seed     string                 `json:"seed"`
	Applied  []*CacheEntry          `json:"applied"`
	Outdated []*CacheEntry          `json:"outdated"`
	nodes    map[string]*CacheEntry `json:"-"`
}

// LoadCacheFromFile loads Cache from a file.
func LoadCacheFromFile(filename string) (*Cache, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &Cache{
				Seed:  data.Timestamp(time.Now()).String(),
				nodes: map[string]*CacheEntry{},
			}, nil
		}
		return nil, err
	}
	return LoadCache(body)
}

// LoadCache loads a cache from in-memory data.
func LoadCache(body []byte) (*Cache, error) {
	c := &Cache{}
	if err := json.Unmarshal(body, c); err != nil {
		return nil, err
	}
	c.nodes = make(map[string]*CacheEntry, len(c.Applied))
	for _, e := range c.Applied {
		if e.NodeName == "" {
			continue
		}
		c.nodes[e.NodeName] = e
	}
	return c, nil
}

// SaveCacheToFile saves Cache to a file.
func SaveCacheToFile(cache *Cache, filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("cannot create a cache file '%v': %v", filename, err)
	}
	defer f.Close()
	if err := cache.Save(f); err != nil {
		return fmt.Errorf("cannot save a cache file '%v': %v", filename, err)
	}
	return nil
}

// CacheFilename computes the filename from a cache entry.
func CacheFilename(ent *CacheEntry) string {
	return CacheFilenameWithName(ent, ent.NodeName, "jsonl")
}

func StateCacheFilename(ent *CacheEntry, stateName string) string {
	return CacheFilenameWithName(ent, stateName, "state")
}

// CacheFilenameWithNodeName computes the filename from a cache entry with a
// custom node name.
func CacheFilenameWithName(ent *CacheEntry, name, ext string) string {
	return fmt.Sprintf("%v-%v-%v.%v", name, ent.Timestamp.Format("20060102150405"), ent.Hash[:8], ext)
}

// Add adds a CacheEntry.
func (c *Cache) Add(e *CacheEntry) {
	c.Applied = append(c.Applied, e)
	if e.NodeName != "" {
		c.nodes[e.NodeName] = e
	}
}

// Get returns i-th CacheEntry in the cache. It returns nil if the index is
// out of bounds.
func (c *Cache) Get(i int) *CacheEntry {
	if i < 0 || i >= len(c.Applied) {
		return nil
	}
	return c.Applied[i]
}

// TruncateAfter removes all cache information after i-th position. For
// example, if the cache has three entries [a, b, c], TruncateAfter(0) removes
// b and c, and only a remains.
func (c *Cache) TruncateAfter(i int) {
	if i >= len(c.Applied)-1 {
		return
	}
	if i < 0 {
		i = -1
	}
	c.Outdated = append(c.Outdated, c.Applied[i+1:]...)
	c.Applied = c.Applied[:i+1]
}

// Node returns a CacheEntry of a data source node (a source or a stream)
// having the name.
func (c *Cache) Node(name string) *CacheEntry {
	return c.nodes[name]
}

// Save saves the Cache.
func (c *Cache) Save(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}
