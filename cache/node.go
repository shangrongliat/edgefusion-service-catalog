package cache

import (
	"errors"
	"sync"

	"edgefusion-service-catalog/model"
	"edgefusion-service-catalog/util"
)

// Cache 结构体用于管理用户
type Cache struct {
	mu          sync.RWMutex              // 读写锁，保护用户数据的并发访问
	local_cache *model.Catalog            // 用户数据存储在内存中的映射
	ecache      map[string]*model.Catalog // 用户数据存储在内存中的映射
	ParentId    string                    //父节点
	Version     string                    // 信息版本
}

// NewCacheManager 创建一个新的 内存管理
func NewCacheManager(localCC *model.Catalog, parentId, version string) *Cache {
	return &Cache{
		local_cache: localCC,
		ecache:      make(map[string]*model.Catalog),
		ParentId:    parentId,
		Version:     version,
	}
}

func (um *Cache) UpdateVersion() {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.Version = util.ToStringUuid()
}

// AddCache 添加修改本地节点缓存
func (um *Cache) AddCache(cache *model.Catalog) {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.local_cache = cache
}

// GetCache 获取本地节点缓存
func (um *Cache) GetCache() *model.Catalog {
	um.mu.RLock()
	defer um.mu.RUnlock()
	return um.local_cache
}

// AddECache 添加其他节点缓存
func (um *Cache) AddECache(cache *model.Catalog) {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.ecache[cache.ID] = cache
}

// GetECache 获取其他节点缓存
func (um *Cache) GetECache(id string) (*model.Catalog, error) {
	um.mu.RLock()
	defer um.mu.RUnlock()
	cache, exists := um.ecache[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return cache, nil
}

// UpdateECache 更新其他节点缓存信息
func (um *Cache) UpdateECache(id string, cache *model.Catalog) error {
	um.mu.Lock()
	defer um.mu.Unlock()
	if _, exists := um.ecache[id]; !exists {
		return errors.New("user not found")
	}
	um.ecache[id] = cache
	return nil
}

// DeleteECache 删除其他节点缓存信息
func (um *Cache) DeleteECache(id string) error {
	um.mu.Lock()
	defer um.mu.Unlock()
	if _, exists := um.ecache[id]; !exists {
		return errors.New("user not found")
	}
	delete(um.ecache, id)
	return nil
}

// ListECache 列出所有外部节点缓存信息
func (um *Cache) ListECache() []*model.Catalog {
	um.mu.RLock()
	defer um.mu.RUnlock()
	caches := make([]*model.Catalog, 0, len(um.ecache))
	for _, user := range um.ecache {
		caches = append(caches, user)
	}
	return caches
}
