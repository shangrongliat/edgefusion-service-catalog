package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"edgefusion-service-catalog/util"
	"github.com/robfig/cron"

	"edgefusion-service-catalog/model"
)

// Cache 结构体用于管理用户
type Cache struct {
	sync.RWMutex                             // 读写锁，保护用户数据的并发访问
	local_cache  *model.Catalog              // 本地节点数据存储在内存中的映射
	ecache       map[string]*model.Catalog   // 外部节点数据存储在内存中的映射
	parentId     string                      // 父节点
	enode        map[string]*model.NodeCache // 父节点下的边计算节点
	version      string                      // 信息版本
	path         string                      // 持久化文件路径
	cron         *cron.Cron                  // 定时器
	StatusChan   chan string
}

// NewCacheManager 创建一个新的 内存管理
func NewCacheManager(path string) *Cache {
	return &Cache{
		local_cache: &model.Catalog{},
		ecache:      make(map[string]*model.Catalog),
		enode:       make(map[string]*model.NodeCache),
		path:        path,
		cron:        cron.New(),
		StatusChan:  make(chan string), // 状态通知通道
	}
}

// AddNodeCache 添加修改本地节点信息缓存与隶属边节点边计算节点对应关系
// 节点的活跃与否只会影响该节点的权重信息，不会影响本地信息的版本
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
				nodeCache.Passing = 50
			} else {
				nodeCache.Passing = 0
			}
		}
	}
	// 修改本地节点缓存时，同时修改版本号
	//c.version = util.ToStringUuid()
}

func (c *Cache) AddServiceCache(ser []*model.Service) {
	c.Lock()
	defer c.Unlock()
	// 更新标记，默认为false，当服务信息发生变化时，进行服务版本更新
	updateFlag := false
	sc := c.local_cache.SC
	for _, service := range ser {
		// 服务ID由节点ID与服务名称进行拼接
		serID := fmt.Sprintf("%s*%s", c.local_cache.ID, service.Name)
		oldCatalog, exist := sc[serID]
		if !exist {
			// 如果本地缓存没有该服务，则添加该服务，并且更新信息版本
			sc[serID] = &model.ServiceCatalog{
				ID:             serID,
				Name:           service.Name,
				Status:         service.Status,
				CheckInterface: "", // 服务状态检测接口，当前为空
				CheckInterval:  "",
				Port:           "", //服务端口
			}
			updateFlag = true
		} else {
			// 如果对应服务存在，则判断其服务状态是否有变化
			if oldCatalog.Status != service.Status {
				// 从新修改状态
				oldCatalog.Status = service.Status
				updateFlag = true
			}
		}
	}
	// 根据修改标记进行版本升级控制
	if updateFlag {
		// 本地服务不在本地consul中注册
		c.version = util.ToStringUuid()
	}
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

// GetCacheBinary 获取二进制缓存信息
func (c *Cache) GetCacheBinary(data []byte) []byte {
	c.RLock()
	defer c.RUnlock()
	var cache model.Broadcast
	// 解析判断，如果解析失败或者请求ID和本机ID不相同，则返回nil,此次询问不做回答
	if err := json.Unmarshal(data, &cache); err != nil {
		log.Printf("数据解析失败,缓存添加.数据: %s . 异常: %v \n", string(data), err)
		return nil
	}
	if cache.ID == c.local_cache.ID {
		log.Printf("询问自循环跳过处理.请求节点ID : %s . 本地节点ID: %s \n", cache.ID, c.local_cache.ID)
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
	// 判断该节点是否在本地存在
	if ecache, ok := c.ecache[cache.ID]; ok {
		// 如果存在，则比对对应的服务状态变化情况
		for sid, sinfo := range cache.SC {
			ecache.SC[sid] = sinfo
		}
	}
	c.ecache[cache.ID] = &cache
	// TODO 将新增的服务ID写入管道中
	c.StatusChan <- cache.ID
}

// GetECache 获取其他节点缓存
func (c *Cache) GetECache(broadcast *model.Broadcast) (*model.Catalog, bool) {
	c.RLock()
	defer c.RUnlock()
	cache, exists := c.ecache[broadcast.ID]
	// 如果缓存中不存在或者版本不相同，则返回nil对象，并且返回false
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

// Load 服务启动加载本地文件缓存
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

// Save 将缓存写入本地文件
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

// StartPersisting 定时本地持久化
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

func (c *Cache) GetLocalCache() *model.Catalog {
	return c.local_cache
}

func (c *Cache) GetLocalSc() map[string]*model.ServiceCatalog {
	return c.local_cache.SC
}

func (c *Cache) GetCacheById(cacheId string) *model.Catalog {
	catalog, ok := c.ecache[cacheId]
	if !ok {
		return nil
	}
	return catalog
}
