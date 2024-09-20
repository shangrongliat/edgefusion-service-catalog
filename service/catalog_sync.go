package service

import (
	"edgefusion-service-catalog/cache"
	"edgefusion-service-catalog/httplink"
	"edgefusion-service-catalog/model"
	"edgefusion-service-catalog/model/consul"
)

type Sync struct {
	cache *cache.Cache
}

func NewCatalogSync(cache *cache.Cache) *Sync {
	return &Sync{
		cache: cache,
	}
}

func (s *Sync) ServiceRegistry() {
	for {
		select {
		// 缓存ID（节点ID），根据服务ID获取对应的信息，将最新的信息注册到consul中
		case id, ok := <-s.cache.StatusChan:
			if !ok {
				return
			}
			// 根据节点ID获取缓存信息，将信息注册到服务中
			catalogCache := s.cache.GetCacheById(id)
			isHighWeight := s.cache.IsBrotherNode(id)
			registry := s.serviceConvertRegistry(catalogCache, isHighWeight)
			for _, register := range registry {
				//拼接HTTP PUT结构体，后调用接口发送到consul中
				httplink.Put(register)
			}
		}
	}
}

func (s *Sync) serviceConvertRegistry(ser *model.Catalog, isHighWeight bool) []*consul.Register {
	var list []*consul.Register
	for _, value := range ser.SC {
		register := &consul.Register{
			ID:                value.ID,
			Name:              value.Name,
			Address:           ser.IP,
			Port:              0,
			Meta:              make(map[string]string),
			EnableTagOverride: false,
			Weights: consul.Weights{
				Passing: 10,
				Warning: 0,
			},
		}
		if isHighWeight {
			register.Weights = consul.Weights{
				Passing: 50,
				Warning: 0,
			}
		}
		if value.Status == "Failed" {
			register.Weights = consul.Weights{
				Passing: 1,
				Warning: 0,
			}
		}
		list = append(list, register)
	}
	return list
}
