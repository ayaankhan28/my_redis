// package main
//
// import (
// 	"fmt"
// 	"net/http"
// 	"sync"
//     "runtime"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/groupcache/lru"
// )
//
// // Custom LRU Cache struct
// type LRUCache struct {
// 	cache *lru.Cache
// 	mu    sync.Mutex
// }
//
// // Create new LRU Cache
// func NewLRUCache(size int) *LRUCache {
// 	return &LRUCache{
// 		cache: lru.New(size),
// 	}
// }
//
// // PUT: Store a key-value pair
// func (c *LRUCache) Put(key string, value string) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	c.cache.Add(key, value)
// }
//
// // GET: Retrieve a value by key
// func (c *LRUCache) Get(key string) (string, bool) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	if value, ok := c.cache.Get(key); ok {
// 		return value.(string), true
// 	}
// 	return "", false
// }
//
// func main() {
//     runtime.GOMAXPROCS(runtime.NumCPU()) // Use all available CPU cores
//
// 	r := gin.Default()
// 	cache := NewLRUCache(100) // Set LRU cache capacity
//
// 	// PUT API
// 	r.POST("/put", func(c *gin.Context) {
// 		var req struct {
// 			Key   string `json:"key"`
// 			Value string `json:"value"`
// 		}
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 			return
// 		}
//
// 		// Store in LRU cache
// 		cache.Put(req.Key, req.Value)
// 		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Key inserted/updated successfully"})
// 	})
//
// 	// GET API
// 	r.GET("/get", func(c *gin.Context) {
// 		key := c.Query("key")
// 		if value, found := cache.Get(key); found {
// 			c.JSON(http.StatusOK, gin.H{"status": "OK", "key": key, "value": value})
// 		} else {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
// 		}
// 	})
//
// 	// Start server
// 	fmt.Println("🚀 Server running on http://localhost:7171")
// 	r.Run(":7171") // Run server on port 7171
// }
package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
)

// Custom LRU Cache struct
type LRUCache struct {
	cache *lru.Cache[string, string]
}

// Create new LRU Cache (thread-safe)
func NewLRUCache(size int) *LRUCache {
	cache, _ := lru.New[string, string](size)
	return &LRUCache{cache: cache}
}

// PUT: Store a key-value pair
func (c *LRUCache) Put(key string, value string) {
	c.cache.Add(key, value)
}

// GET: Retrieve a value by key
func (c *LRUCache) Get(key string) (string, bool) {
	if value, found := c.cache.Get(key); found {
		return value, true
	}
	return "", false
}

func main() {
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Setup Gin router
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	// Create LRU Cache
	cache := NewLRUCache(100000) // Increased cache size

	// PUT API
	r.POST("/put", func(c *gin.Context) {
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		cache.Put(req.Key, req.Value)
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Key inserted/updated successfully"})
	})

	// GET API
	r.GET("/get", func(c *gin.Context) {
		key := c.Query("key")
		if value, found := cache.Get(key); found {
			c.JSON(http.StatusOK, gin.H{"status": "OK", "key": key, "value": value})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		}
	})

	// HTTP Server with Keep-Alive and timeouts
	srv := &http.Server{
		Addr:         ":7171",
		Handler:      r,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	fmt.Println("🚀 Server running on http://localhost:7171")
	log.Fatal(srv.ListenAndServe())
}
