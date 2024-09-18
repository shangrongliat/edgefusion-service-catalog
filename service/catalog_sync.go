package service

import (
	"fmt"

	"edgefusion-service-catalog/cache"
)

type Sync struct {
	cache *cache.Cache
}

func NewCatalogSync(cache *cache.Cache) *Sync {
	return &Sync{
		cache: cache,
	}
}

func (s *Sync) serviceRegistry() {
	go func() {
		for {
			select {
			// 缓存ID（节点ID），根据服务ID获取对应的信息，将最新的信息注册到consul中
			case id, ok := <-s.cache.StatusChan:
				if !ok {
					return
				}
				// 根据节点ID获取缓存信息，将信息注册到服务中
				catalogCache := s.cache.GetCacheById(id)
				//拼接HTTP PUT结构体，后调用接口发送到consul中
				fmt.Println(catalogCache)
			}
		}
	}()
}
