package cache

import (
	"encoding/json"
	"errors"
	"github.com/robfig/cron"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"edgefusion-service-catalog/model"
)

// Cache 结构体用于管理用户
type Cache struct {
	sync.RWMutex                             // 读写锁，保护用户数据的并发访问
	local_cache  *model.Catalog              // 用户数据存储在内存中的映射
	ecache       map[string]*model.Catalog   // 用户数据存储在内存中的映射
	parentId     string                      //父节点
	enode        map[string]*model.NodeCache //父节点下的边计算节点
	version      string                      // 信息版本
	path         string                      // 持久化文件路径
	cron         *cron.Cron                  // 定时器
}

// NewCacheManager 创建一个新的 内存管理
func NewCacheManager(path, parentId string) *Cache {
	return &Cache{
		local_cache: &model.Catalog{},
		ecache:      make(map[string]*model.Catalog),
		enode:       make(map[string]*model.NodeCache),
		parentId:    parentId,
		path:        path,
		cron:        cron.New(),
	}
}

// AddNodeCache 添加修改本地节点信息缓存与隶属边节点边计算节点对应关系
func (c *Cache) AddNodeCache(cache *model.Node) {
	c.Lock()
	defer c.Unlock()
	// 更新基本信息
	c.local_cache.ID = cache.Node.ID
	c.local_cache.IP = cache.Node.IP
	c.parentId = cache.Parent.ID
	// 更新管理节点与边计算节点关系
	for _, enode := range cache.ENodeList {
		nodeCache, ok := c.enode[enode.ID]
		if !ok && enode.State == "Activated" {
			// 缓存中不存在此边界点，并且节点状态为活跃，则新增边计算节点
			c.enode[enode.ID] = &model.NodeCache{
				ID:      enode.IP,
				Name:    enode.Name,
				State:   enode.State,
				Passing: 50,
				Warning: 0,
			}
		} else if ok && enode.State != nodeCache.State {
			// 缓存中存在该边计算节点，并且状态与新的信息状态不同
			// 先更新缓存状态
			nodeCache.State = enode.State
			// 根据新信息修改权重
			if enode.State == "Activated" {
				nodeCache.Warning = 0
			} else {
				nodeCache.Warning = 1
			}
		}
	}
	// 修改本地节点缓存时，同时修改版本号
	//c.version = util.ToStringUuid()
}

func (c *Cache) AddServiceCache(ser *model.Service) {

}

// GetCache 获取本地节点缓存
func (c *Cache) GetCache() *model.Catalog {
	c.RLock()
	defer c.RUnlock()
	return c.local_cache
}

func (c *Cache) GetBroadcastInfo() (br model.Broadcast) {
	c.RLock()
	defer c.RUnlock()
	br.ID = c.local_cache.ID
	br.Version = c.version
	return
}

func (c *Cache) GetCacheBinary(data []byte) []byte {
	c.RLock()
	defer c.RUnlock()
	var cache model.Broadcast
	// 解析判断，如果解析失败或者请求ID和本机ID不相同，则返回nil,此次询问不做回答
	if err := json.Unmarshal(data, &cache); err != nil || cache.ID != c.local_cache.ID {
		log.Printf("数据解析失败,缓存添加.数据: %s . 异常: %v \n", string(data), err)
		return nil
	}
	marshal, err := json.Marshal(c.local_cache)
	if err != nil {
		log.Printf("本地数据json化失败.%v \n", err)
	}
	return marshal
}

// AddECache 添加其他节点缓存
func (c *Cache) AddECache(cache *model.Catalog) {
	c.Lock()
	defer c.Unlock()
	c.ecache[cache.ID] = cache
}

// AddECache 添加其他节点缓存
func (c *Cache) AddECacheBinary(data []byte) {
	c.Lock()
	defer c.Unlock()
	var cache model.Catalog
	if err := json.Unmarshal(data, &cache); err != nil {
		log.Printf("数据解析失败,缓存添加.数据: %s . 异常: %v \n", string(data), err)
		return
	}
	c.ecache[cache.ID] = &cache

}

// GetECache 获取其他节点缓存
func (c *Cache) GetECache(broadcast *model.Broadcast) (*model.Catalog, bool) {
	c.RLock()
	defer c.RUnlock()
	cache, exists := c.ecache[broadcast.ID]
	//
	if !exists || cache.Version != broadcast.Version {
		return nil, false
	}
	return cache, true
}

// UpdateECache 更新其他节点缓存信息
func (c *Cache) UpdateECache(id string, cache *model.Catalog) error {
	c.Lock()
	defer c.Unlock()
	if _, exists := c.ecache[id]; !exists {
		return errors.New("user not found")
	}
	c.ecache[id] = cache
	return nil
}

// DeleteECache 删除其他节点缓存信息
func (c *Cache) DeleteECache(id string) error {
	c.Lock()
	defer c.Unlock()
	if _, exists := c.ecache[id]; !exists {
		return errors.New("user not found")
	}
	delete(c.ecache, id)
	return nil
}

// ListECache 列出所有外部节点缓存信息
func (c *Cache) ListECache() []*model.Catalog {
	c.RLock()
	defer c.RUnlock()
	caches := make([]*model.Catalog, 0, len(c.ecache))
	for _, user := range c.ecache {
		caches = append(caches, user)
	}
	return caches
}

func (c *Cache) Load() error {
	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，缓存为空
		}
		return err
	}
	var cacheData model.CacheData
	if err := json.Unmarshal(data, &cacheData); err != nil {
		return err
	}
	c.Lock()
	defer c.Unlock()
	c.local_cache = cacheData.LocalCache
	c.ecache = cacheData.Ecache
	c.parentId = cacheData.ParentId
	c.version = cacheData.Version
	return nil
}

func (c *Cache) Save() error {
	c.Lock()
	defer c.Unlock()
	data, err := json.Marshal(model.CacheData{
		LocalCache: c.local_cache,
		Ecache:     c.ecache,
		ParentId:   c.parentId,
		Version:    c.version,
	})
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(c.path, data, 0644); err != nil {
		return err
	}
	return nil
}

func (c *Cache) StartPersisting(interval time.Duration) {
	job := func() {
		if err := c.Save(); err != nil {
			log.Println("Error saving cache:", err)
		}
	}
	if err := c.cron.AddFunc("@every 60s", job); err != nil {
		return
	}
	c.cron.Start()
}
