package cache

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu       sync.RWMutex
	entries  map[string]cacheEntry // is stored locally at the end of session
	temp     map[string]cacheEntry // is not stored
	interval time.Duration
}

type LocalCacheStorage struct {
	Entries map[string]SerializableCacheEntry `json:"entries"`
}

type SerializableCacheEntry struct {
	CreatedAt time.Time `json:"createdAt"`
	Data      string    `json:"data"` // Base64 encoded string
}

func NewCache() *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		temp:     make(map[string]cacheEntry),
		interval: time.Duration(time.Hour * 720), // Making 1 month interval
	}
	return c
}

func NewCacheFromFile() *Cache {
	entries, err := loadCacheFile()
	if err != nil {
		log.Println("--DEBUG-- error with loadCacheFile(), using a empty cache entries map")
	} else {
		log.Println("--DEBUG-- NewCacheFromFile() OK")
	}
	c := &Cache{
		entries:  entries,
		temp:     make(map[string]cacheEntry),
		interval: time.Duration(time.Hour * 720), // Making 1 month interval
	}
	return c
}

func (c *Cache) Add(key string, val []byte) {
	// Handle mutex
	c.mu.Lock()
	defer c.mu.Unlock()

	// Add to cache
	c.entries[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val:       val,
	}
}

func (c *Cache) AddToTemp(key string, val []byte) {
	// Handle mutex
	c.mu.Lock()
	defer c.mu.Unlock()

	// Add to temp cache
	c.temp[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	// Handle mutex
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, found := c.entries[key]
	// Check if cache entry exists and if it's not too old
	if !found || time.Since(entry.createdAt) > c.interval {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) GetFromTemp(key string) ([]byte, bool) {
	// Handle mutex
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, found := c.temp[key]
	// Check if cache entry exists and if it's not too old
	if !found || time.Since(entry.createdAt) > c.interval {
		return nil, false
	}
	return entry.val, true
}

func getLocalCacheStoragePath() string {
	outputDir := ""

	// Check if OS used is Windows or else
	if runtime.GOOS == "windows" {
		// For windows builds, Get the path to the currently running executable
		execPath, err := os.Executable()
		if err != nil {
			log.Fatalf("Error with getLocalCacheStoragePath(): %v", err)
		} else {
			outputDir = filepath.Dir(execPath)
		}
	} else {
		// For Linux/Mac, use working directory + client directory
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error with getLocalCacheStoragePath(): %v", err)
		} else {
			outputDir = filepath.Join(workingDir, "client")
		}
	}

	// If we couldn't determine a directory, fall back to current directory
	if outputDir == "" {
		outputDir = "."
	}

	// Build the path to the "config/cache.json" file
	return filepath.Join(outputDir, "config", "cache.json")
}

func loadCacheFile() (map[string]cacheEntry, error) {
	emptyMap := make(map[string]cacheEntry)

	// Open local cache.json file
	localCacheFilePath := getLocalCacheStoragePath()
	f, err := os.Open(localCacheFilePath)
	if err != nil {
		log.Printf("--ERROR-- with loadCacheFile(), couldn't open cache.json: %v\n", err)
		return emptyMap, err
	}
	defer f.Close()

	// Define a temporary structure to hold the serialized data
	serializedCache := LocalCacheStorage{}

	// Read data from file
	err = json.NewDecoder(f).Decode(&serializedCache)
	if err != nil {
		log.Printf("--ERROR-- with localCacheFile(), couldn't decode cache.json: %v\n", err)
		return emptyMap, err
	}

	// Convert the serialized entries back to cacheEntry format
	result := make(map[string]cacheEntry)

	for key, entry := range serializedCache.Entries {
		// Decode the Base64 string back to binary
		decodedData, err := base64.StdEncoding.DecodeString(entry.Data)
		if err != nil {
			log.Printf("--WARNING-- with localCacheFile(), couldn't decode Base64 data for key %s: %v\n", key, err)
			continue
		}

		result[key] = cacheEntry{
			createdAt: entry.CreatedAt,
			val:       decodedData,
		}
	}

	// Return data
	log.Println("--DEBUG-- loadCacheFile() OK")
	return result, nil
}

func (c *Cache) DumpCacheFile() {
	// Create/erase local cache.json file
	localCacheFilePath := getLocalCacheStoragePath()
	f, err := os.Create(localCacheFilePath)
	if err != nil {
		log.Printf("--ERROR-- with DumpCacheFile, couldn't create cache.json: %v\n", err)
		return
	}
	defer f.Close()

	// Convert cache entries to serializable format
	serializableEntries := make(map[string]SerializableCacheEntry)

	for key, entry := range c.entries {
		// Convert the binary data to Base64
		base64Data := base64.StdEncoding.EncodeToString(entry.val)

		serializableEntries[key] = SerializableCacheEntry{
			CreatedAt: entry.createdAt,
			Data:      base64Data,
		}
	}

	// Create the structure to marchal
	localCache := LocalCacheStorage{
		Entries: serializableEntries,
	}

	// Marshal and write
	data, err := json.Marshal(localCache)
	if err != nil {
		log.Printf("--ERROR-- with DumpCacheFile, couldn't json.Marshal data: %v\n", err)
		return
	}

	// Write to file
	_, err = f.Write(data)
	if err != nil {
		log.Printf("--ERROR-- with DumpCacheFile, couldn't write data in cache.json: %v\n", err)
		return
	}
}
